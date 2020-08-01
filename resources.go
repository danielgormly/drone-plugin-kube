package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	appV1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	coreV1 "k8s.io/api/core/v1"
	extV1BetaV1 "k8s.io/api/extensions/v1beta1"
	netV1BetaV1 "k8s.io/api/networking/v1beta1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment -- Checks if deployment already exists, updates if it does, creates if it doesn't
func CreateOrUpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appV1.Deployment, hpa *autoscalingV1.HorizontalPodAutoscaler) error {
	deploymentExists, err := deploymentExists(clientset, namespace, deployment.Name)
	if err != nil {
		return err
	}
	if deploymentExists {
		log.Printf("📦 Found existing deployment '%s'. Updating.", deployment.Name)
		_, err = clientset.AppsV1().Deployments(namespace).Update(deployment)
		return err
	}
	log.Printf("📦 Creating new deployment '%s'. Updating.", deployment.Name)
	_, err = clientset.AppsV1().Deployments(namespace).Create(deployment)
	if err != nil {
		return err
	}

	if hpa != nil {
		return ApplyHorizontalAutoscaler(clientset, namespace, hpa)
	}

	return nil
}

// deploymentExists -- Updates given deployment in Kubernetes
func deploymentExists(clientset *kubernetes.Clientset, namespace string, name string) (bool, error) {
	_, err := clientset.AppsV1().Deployments(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		// TODO: Only conver to StatusError if the error is in fact a status error
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

const (
	stateUpdated = "📦 Updated"
	stateFailed  = "⛔️ Failed"
)

// waitUntilDeploymentSettles -- Waits until ready, failure or timeout
func waitUntilDeploymentSettles(clientset *kubernetes.Clientset, namespace string, name string, timeoutInSeconds int64) (string, error) {
	fieldSelector := strings.Join([]string{"metadata.name", name}, "=")
	watchOptions := metaV1.ListOptions{
		FieldSelector: fieldSelector,
		Watch:         true,
	}

	watcher, err := clientset.AppsV1().Deployments(namespace).Watch(watchOptions)
	if err != nil {
		return stateFailed, fmt.Errorf("watch deployment; %w", err)
	}

	liveDeployment, err := clientset.AppsV1().Deployments(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return stateFailed, fmt.Errorf("get deployment; %w", err)
	}

	log.Printf("📦 Unavailable replicas: %d", liveDeployment.Status.UnavailableReplicas)
	if liveDeployment.Status.UnavailableReplicas == 0 {
		return "📦 Updated", nil
	}

	timer := time.NewTimer(time.Duration(timeoutInSeconds) * time.Second)
	for {
		select {
		case <-timer.C:
			return stateFailed, errors.New("deployment watcher timed out. Something is wrong")
		case event := <-watcher.ResultChan():
			deployment := event.Object.(*appV1.Deployment)
			if deployment.Status.UnavailableReplicas == 0 {
				return stateUpdated, nil
			}
			log.Printf("📦 Unavailable replicas: %d", deployment.Status.UnavailableReplicas)
		}
	}
}

// ApplyService creates a service if it doesn't exists, updates it if it does
func ApplyService(clientset *kubernetes.Clientset, namespace string, service *coreV1.Service) error {
	existingService, exists, err := getService(clientset, namespace, service.Name)
	if err != nil {
		return err
	}

	if exists {
		_, err = clientset.CoreV1().Services(namespace).Update(existingService)
		return err
	}

	_, err = clientset.CoreV1().Services(namespace).Create(service)
	return err
}

// getService returns a service object, or false if it doesn't exist
func getService(clientset *kubernetes.Clientset, namespace string, name string) (*coreV1.Service, bool, error) {
	service, err := clientset.CoreV1().Services(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return service, true, nil
}

// ApplyConfigMap -- Updates given deployment in Kubernetes
func ApplyConfigMap(clientset *kubernetes.Clientset, namespace string, configMap *coreV1.ConfigMap, configMapPath string) error {
	log.Printf("📦 Reading contents of %s", configMapPath)
	info, err := os.Stat(configMapPath)
	if err != nil {
		return err
	}

	configMapData := make(map[string]string)
	mode := info.Mode()

	if mode.IsDir() {
		err = filepath.Walk(configMapPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() { // Don't process directories, let `Walk` walk its files
				return nil
			}

			fileContents, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(configMapPath, path)
			if err != nil {
			    return err
			}

			// Replace slashes with dashes because kube doesn't like em
			configMapData[strings.ReplaceAll(relativePath, "/", "-")] = string(fileContents)
			return nil
		})

		if err != nil {
			return err
		}
	} else if mode.IsRegular() {
		_, filename := filepath.Split(configMapPath)

		fileContents, err := ioutil.ReadFile(configMapPath)
		if err != nil {
			return err
		}

		configMapData[filename] = string(fileContents)
	}

	configMap.Data = configMapData

	// Check if deployment exists
	exists, err := configMapExists(clientset, namespace, configMap.Name)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("📦 Found existing deployment. Updating %s.", configMap.Name)
		_, err = clientset.CoreV1().ConfigMaps(namespace).Update(configMap)
		return err
	}

	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(configMap)
	return err
}

// configMapExists -- Updates given deployment in Kubernetes
func configMapExists(clientset *kubernetes.Clientset, namespace string, name string) (bool, error) {
	_, err := clientset.CoreV1().ConfigMaps(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ApplyExtensionsV1beta1Ingress(clientset *kubernetes.Clientset, namespace string, ingress *extV1BetaV1.Ingress) error {
	_, exists, err := getExtensionsV1beta1Ingress(clientset, namespace, ingress.Name)
	if err != nil {
		return err
	}

	if !exists {
		_, err = clientset.ExtensionsV1beta1().Ingresses(namespace).Create(ingress)
		return err
	}

	_, err = clientset.ExtensionsV1beta1().Ingresses(namespace).Update(ingress)
	return err
}

func getExtensionsV1beta1Ingress(clientset *kubernetes.Clientset, namespace string, name string) (*extV1BetaV1.Ingress, bool, error) {
	ingress, err := clientset.ExtensionsV1beta1().Ingresses(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return ingress, true, nil
}

func ApplyNetworkingV1beta1Ingress(clientset *kubernetes.Clientset, namespace string, ingress *netV1BetaV1.Ingress, additionalAnnotations map[string]string) error {
	_, exists, err := getNetworkingV1beta1Ingress(clientset, namespace, ingress.Name)
	if err != nil {
		return err
	}

	for k, v := range additionalAnnotations {
		ingress.Annotations[k] = v
	}

	if !exists {
		_, err = clientset.NetworkingV1beta1().Ingresses(namespace).Create(ingress)
		return err
	}

	_, err = clientset.NetworkingV1beta1().Ingresses(namespace).Update(ingress)
	return err
}

func getNetworkingV1beta1Ingress(clientset *kubernetes.Clientset, namespace string, name string) (*netV1BetaV1.Ingress, bool, error) {
	ingress, err := clientset.NetworkingV1beta1().Ingresses(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return ingress, true, nil
}

func ApplySecret(clientset *kubernetes.Clientset, namespace string, secret *coreV1.Secret, data map[string]string) error {
	for k, v := range data {
		if secret.StringData == nil {
			secret.StringData = make(map[string]string)
		}
		secret.StringData[k] = v
	}

	_, exists, err := getSecret(clientset, namespace, secret.Name)
	if !exists {
		_, err := clientset.CoreV1().Secrets(namespace).Create(secret)
		return err
	}

	_, err = clientset.CoreV1().Secrets(namespace).Update(secret)
	return err
}

func getSecret(clientset *kubernetes.Clientset, namespace, name string) (*coreV1.Secret, bool, error) {
	secret, err := clientset.CoreV1().Secrets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	return secret, true, nil
}

func ApplyHorizontalAutoscaler(clientset *kubernetes.Clientset, namespace string, autoscaler *autoscalingV1.HorizontalPodAutoscaler) error {
	_, exists, err := getAutoscaler(clientset, namespace, autoscaler.Name)
	if err != nil {
		return err
	}

	if !exists {
		_, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Create(autoscaler)
		return err
	}

	_, err = clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Update(autoscaler)
	return err
}

func getAutoscaler(clientset *kubernetes.Clientset, namespace, name string)  (*autoscalingV1.HorizontalPodAutoscaler,bool,  error) {
	as, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		statusError, ok := err.(*kubeErrors.StatusError)
		if ok && statusError.Status().Code == http.StatusNotFound {
			return nil, false, nil
		}

		return nil, false, err
	}

	return as, true, nil
}
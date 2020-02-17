package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	v1BetaV1 "k8s.io/api/extensions/v1beta1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment -- Checks if deployment already exists, updates if it does, creates if it doesn't
func CreateOrUpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appV1.Deployment) error {
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
	return err
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

// waitUntilDeploymentSettled -- Waits until ready, failure or timeout
func waitUntilDeploymentSettled(clientset *kubernetes.Clientset, namespace string, name string, timeoutInSeconds int64) (string, error) {
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

// ApplyConfigMapFromFile -- Updates given deployment in Kubernetes
func ApplyConfigMapFromFile(clientset *kubernetes.Clientset, namespace string, configMap *coreV1.ConfigMap, path string) error {
	log.Printf("📦 Reading contents of %s", path)
	_, filename := filepath.Split(path)
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	configMapData := make(map[string]string)
	configMapData[filename] = string(fileContents)
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

func ApplyIngress(clientset *kubernetes.Clientset, namespace string, ingress *v1BetaV1.Ingress) error {
	_, exists, err := getIngress(clientset, namespace, ingress.Name)
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

func getIngress(clientset *kubernetes.Clientset, namespace string, name string) (*v1BetaV1.Ingress, bool, error) {
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

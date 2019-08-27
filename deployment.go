package main

import (
	"log"
	"strings"

	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment -- Checks if deployment already exists, updates if it does, creates if it doesn't
func CreateOrUpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
	deploymentExists, err := deploymentExists(clientset, namespace, deployment.Name)
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
	_, err := clientset.AppsV1().Deployments(namespace).Get(name, meta.GetOptions{})
	if err != nil {
		// TODO: Only conver to StatusError if the error is in fact a status error
		statusError, ok := err.(*errors.StatusError)
		if ok == true && statusError.Status().Code == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// waitUntilDeploymentSettled -- Waits until ready, failure or timeout
func waitUntilDeploymentSettled(clientset *kubernetes.Clientset, namespace string, name string, timeout int64) (state string, err error) {
	fieldSelector := strings.Join([]string{"metadata.name", name}, "=")
	watchOptions := meta.ListOptions{
		FieldSelector: fieldSelector,
		Watch:         true,
	}
	watcher, error := clientset.AppsV1().Deployments(namespace).Watch(watchOptions)
	liveDeployment, error := clientset.AppsV1().Deployments(namespace).Get(name, meta.GetOptions{})
	log.Printf("📦 Unavailable replicas: %d", liveDeployment.Status.UnavailableReplicas)
	if liveDeployment.Status.UnavailableReplicas == 0 {
		return "📦 Updated", error
	}
	i := 0
	for {
		event := <-watcher.ResultChan()
		deployment := event.Object.(*appv1.Deployment)
		if deployment.Status.UnavailableReplicas == 0 {
			return "📦 Updated", error
		}
		log.Printf("📦 Unavailable replicas: %d", deployment.Status.UnavailableReplicas)
		i++
	}
}

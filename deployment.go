package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	appv1 "k8s.io/api/apps/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment -- Checks if deployment already exists, updates if it does, creates if it doesn't
func CreateOrUpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
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
	_, err := clientset.AppsV1().Deployments(namespace).Get(name, meta.GetOptions{})
	if err != nil {
		// TODO: Only conver to StatusError if the error is in fact a status error
		statusError, ok := err.(*k8sErrors.StatusError)
		if ok == true && statusError.Status().Code == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

const stateFailed = "⛔️ Failed"

// waitUntilDeploymentSettled -- Waits until ready, failure or timeout
func waitUntilDeploymentSettled(clientset *kubernetes.Clientset, namespace string, name string, timeoutInSeconds int64) (string, error) {
	fieldSelector := strings.Join([]string{"metadata.name", name}, "=")
	watchOptions := meta.ListOptions{
		FieldSelector: fieldSelector,
		Watch:         true,
	}

	watcher, err := clientset.AppsV1().Deployments(namespace).Watch(watchOptions)
	if err != nil {
		return stateFailed, fmt.Errorf("watch deployment; %w", err)
	}

	liveDeployment, err := clientset.AppsV1().Deployments(namespace).Get(name, meta.GetOptions{})
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
			deployment := event.Object.(*appv1.Deployment)
			if deployment.Status.UnavailableReplicas == 0 {
				return "📦 Updated", nil
			}
			log.Printf("📦 Unavailable replicas: %d", deployment.Status.UnavailableReplicas)
		}
	}
}

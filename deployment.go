package main

import (
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateOrUpdateDeployment -- Checks if deployment already exists, updates if it does, creates if it doesn't
func CreateOrUpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
	deploymentExists, err := deploymentExists(clientset, namespace, deployment.Name)
	if deploymentExists {
		// log.Printf("📦 Found existing deployment. Updating.\n%s\n", depYaml)
		_, err = clientset.AppsV1().Deployments(namespace).Update(deployment)
		return err
	}
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

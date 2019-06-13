package main

import (
	appv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// CreateDeploymentObj -- Construct KubeClient ready json from YAML definition file
func CreateDeploymentObj(yaml string) *appv1.Deployment {
	deployment := appv1.Deployment{}
	scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, &deployment)
	return &deployment
}

// UpdateDeployment -- Updates given deployment in Kubernetes
func UpdateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
	_, err := clientset.AppsV1().Deployments(namespace).Update(deployment)
	return err
}

// CreateDeployment -- Updates given deployment in Kubernetes
func CreateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
	_, err := clientset.AppsV1().Deployments(namespace).Create(deployment)
	return err
}

// DeploymentExists -- Updates given deployment in Kubernetes
func DeploymentExists(clientset *kubernetes.Clientset, namespace string, deploymentName string) (bool, error) {
	_, err := clientset.AppsV1().Deployments(namespace).Get(deploymentName, meta.GetOptions{})
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

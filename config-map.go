package main

import (
	appv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// CreateConfigMap -- Updates given deployment in Kubernetes
func CreateConfigMapFromFile(clientset *kubernetes.Clientset, namespace string, name string) error {
	configMap := appv1.ConfigMap{}
	_, err := clientset.AppsV1().ConfigMap(namespace).Create(configMap)
	return err
}

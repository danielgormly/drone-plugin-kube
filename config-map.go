package main

import (
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// ApplyConfigMapFromFile -- Updates given deployment in Kubernetes
func ApplyConfigMapFromFile(clientset *kubernetes.Clientset, namespace string, configmap *corev1.ConfigMap, path string) error {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	configMapData := make(map[string][]byte)
	configMapData["notsure"] = fileContents
	configmap.BinaryData = configMapData
	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(configmap)
	return err
}

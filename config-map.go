package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ApplyConfigMapFromFile -- Updates given deployment in Kubernetes
func ApplyConfigMapFromFile(clientset *kubernetes.Clientset, namespace string, configmap *corev1.ConfigMap, path string) error {
	log.Printf("📦 Reading contents of %s", path)
	_, filename := filepath.Split(path)
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	configMapData := make(map[string]string)
	configMapData[filename] = string(fileContents)
	configmap.Data = configMapData
	// Check if deployment exists
	deploymentExists, err := configMapExists(clientset, namespace, configmap.Name)
	if deploymentExists {
		log.Printf("📦 Found existing deployment. Updating %s.", configmap.Name)
		_, err = clientset.CoreV1().ConfigMaps(namespace).Update(configmap)
		return err
	}
	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(configmap)
	return err
}

// configMapExists -- Updates given deployment in Kubernetes
func configMapExists(clientset *kubernetes.Clientset, namespace string, name string) (bool, error) {
	_, err := clientset.CoreV1().ConfigMaps(namespace).Get(name, meta.GetOptions{})
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

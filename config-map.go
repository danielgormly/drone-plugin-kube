package main

import (
	appv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// CreateConfigMapObj -- Construct ConfigMap ready json from YAML definition file
func CreateConfigMapObj(yaml string) *appv1.ConfigMap {
	configMap := appv1.ConfigMap{}
	scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, &configMap)
	return &configMap
}

// CreateDeployment -- Updates given deployment in Kubernetes
func CreateDeployment(clientset *kubernetes.Clientset, namespace string, deployment *appv1.Deployment) error {
	_, err := clientset.AppsV1().Deployments(namespace).Create(deployment)
	return err
}

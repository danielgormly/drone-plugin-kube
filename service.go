package main

import (
	"net/http"

	coreV1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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

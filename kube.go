package main

import (
	"encoding/base64"
	"fmt"
	"log"

	appv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// CreateKubeClient -- Creates KubeClient
func (p Plugin) CreateKubeClient() (*kubernetes.Clientset, error) {
	config := clientcmdapi.NewConfig()
	clusterConfig := clientcmdapi.Cluster{
		Server: p.KubeConfig.Server,
	}
	if p.KubeConfig.InsecureSkipTLSVerify == true {
		clusterConfig.InsecureSkipTLSVerify = true
		log.Println("InsecureSkipTLSVerify flag set")
	} else {
		ca, err := base64.StdEncoding.DecodeString(p.KubeConfig.Ca)
		if err != nil {
			log.Fatal(err)
		}
		clusterConfig.CertificateAuthorityData = ca
	}
	config.Clusters["default"] = &clusterConfig
	config.AuthInfos["default"] = &clientcmdapi.AuthInfo{
		Token: p.KubeConfig.Token,
	}
	config.Contexts["default"] = &clientcmdapi.Context{
		Cluster:   "default",
		AuthInfo:  "default",
		Namespace: p.KubeConfig.Namespace,
	}
	config.CurrentContext = "default"
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*config, "default", &clientcmd.ConfigOverrides{}, nil)
	actualCfg, err := clientBuilder.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		fmt.Printf("Unexpected error: %v", err)
	}
	return kubernetes.NewForConfig(actualCfg)
}

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

// ListDeployments -- List deployments in Kubernetes
// func ListDeployments(clientset *kubernetes.Clientset, namespace string) {
// 	deployments, err := clientset.AppsV1().Deployments(namespace).List(v1.ListOptions{})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	fmt.Println(deployments.Items)
// 	// return deployments.Items
// }

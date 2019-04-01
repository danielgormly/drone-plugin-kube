package main

import (
	"fmt"
	"log"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// CreateKubeClient -- Creates KubeClient
func (p Plugin) CreateKubeClient() (*kubernetes.Clientset, error) {
	// ca, err := base64.StdEncoding.DecodeString(p.KubeConfig.Ca)
	config := api.NewConfig()
	config.Clusters["default"] = &api.Cluster{
		Server: p.KubeConfig.Server,
		// CertificateAuthorityData: ca,
		InsecureSkipTLSVerify: true,
	}
	config.AuthInfos["default"] = &api.AuthInfo{
		Token: p.KubeConfig.Token,
	}
	config.Contexts["default"] = &api.Context{
		Cluster:  "default",
		AuthInfo: "default",
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

// WatchPodCounts -- Example function
func WatchPodCounts(clientset *kubernetes.Clientset) {
	for {
		pods, err := clientset.Core().Pods("").List(v1.ListOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		time.Sleep(10 * time.Second)
	}
}

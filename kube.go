package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
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
			log.Printf("possible corrupted CA, or not base64 encoded: %s\n", err)
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
		return nil, fmt.Errorf("client builder client config; %w", err)
	}
	return kubernetes.NewForConfig(actualCfg)
}

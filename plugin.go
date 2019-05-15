package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/aymerick/raymond"
)

type (
	// KubeConfig -- Contains connection settings for Kube client
	KubeConfig struct {
		Ca                    string
		Server                string
		Token                 string
		Namespace             string
		InsecureSkipTLSVerify bool
	}
	// Plugin -- Contains config for plugin
	Plugin struct {
		Template   string
		KubeConfig KubeConfig
	}
)

// Exec -- Runs plugin
func (p Plugin) Exec() error {
	if p.KubeConfig.Server == "" {
		log.Fatal("PLUGIN_SERVER is not defined")
	}
	if p.KubeConfig.Token == "" {
		log.Fatal("PLUGIN_TOKEN is not defined")
	}
	if p.KubeConfig.Ca == "" {
		log.Fatal("PLUGIN_CA is not defined")
	}
	if p.KubeConfig.Namespace == "" {
		p.KubeConfig.Namespace = "default"
	}
	if p.Template == "" {
		log.Fatal("PLUGIN_TEMPLATE, or template must be defined")
	}
	// Make map of environment variables set by Drone
	ctx := make(map[string]string)
	pluginEnv := os.Environ()
	for _, value := range pluginEnv {
		re := regexp.MustCompile(`^PLUGIN_(.*)=(.*)`)
		if re.MatchString(value) {
			matches := re.FindStringSubmatch(value)
			key := strings.ToLower(matches[1])
			ctx[key] = matches[2]
		}
	}
	// TODO: Remove in first release
	fmt.Printf("%#v", ctx)
	// Grab template from filesystem
	raw, err := ioutil.ReadFile(p.Template)
	if err != nil {
		log.Print("Error reading template file:")
		return err
	}
	// Parse template
	depYaml, err := raymond.Render(string(raw), ctx)
	fmt.Printf("%s", depYaml)
	if err != nil {
		panic(err)
	}
	// Connect to Kubernetes
	clientset, err := p.CreateKubeClient()
	if err != nil {
		log.Fatal(err.Error())
	}
	deployment := CreateDeploymentObj(depYaml)
	UpdateDeployment(clientset, p.KubeConfig.Namespace, deployment)
	return err
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/aymerick/raymond"
)

type (
	KubeConfig struct {
		Ca        string
		Endpoint  string
		Token     string
		Namespace string
		Template  string
	}
	Plugin struct {
		Template   string
		KubeConfig KubeConfig
	}
)

func (p Plugin) Exec() error {
	if p.KubeConfig.Endpoint == "" {
		log.Fatal("PLUGIN_ENDPOINT is not defined")
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
	// // connect to Kubernetes
	// clientset, err := p.createKubeClient()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	raw, err := ioutil.ReadFile(p.Template)
	if err != nil {
		log.Print("Error reading template file:")
		return err
	}

	source := string(raw)

	ctx := make(map[string]string)
	ctx["KUBE_CA"] = p.KubeConfig.Ca
	ctx["KUBE_TOKEN"] = p.KubeConfig.Token
	ctx["KUBE_ENDPOINT"] = p.KubeConfig.Endpoint
	ctx["KUBE_NAMESPACE"] = p.KubeConfig.Namespace
	droneEnv := os.Environ()
	for _, value := range droneEnv {
		re := regexp.MustCompile(`^(DRONE_.*)=(.*)`)
		if re.MatchString(value) {
			matches := re.FindStringSubmatch(value)
			ctx[matches[1]] = matches[2]
		}
	}

	// parse template
	tpl, err := raymond.Parse(source)
	if err != nil {
		panic(err)
	}

	result, err := tpl.Exec(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Print(result)

	return err
}

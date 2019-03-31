package main

import (
	"log"
	"os"
)

func main() {
	plugin := Plugin{
		Template: os.Getenv("PLUGIN_TEMPLATE"),
		KubeConfig: KubeConfig{
			Token:     os.Getenv("PLUGIN_TOKEN"),
			Endpoint:  os.Getenv("PLUGIN_ENDPOINT"),
			Ca:        os.Getenv("PLUGIN_CA"),
			Namespace: os.Getenv("PLUGIN_NAMESPACE"),
		},
	}
	err := plugin.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	plugin := Plugin{
		Template:      os.Getenv("PLUGIN_TEMPLATE"),
		ConfigMapFile: os.Getenv("PLUGIN_CONFIGMAP_FILE"),
		HpaTemplate:   os.Getenv("PLUGIN_HPA_TEMPLATE"),
		KubeConfig: KubeConfig{
			Token:                 os.Getenv("PLUGIN_TOKEN"),
			Server:                os.Getenv("PLUGIN_SERVER"),
			Ca:                    os.Getenv("PLUGIN_CA"),
			Namespace:             os.Getenv("PLUGIN_NAMESPACE"),
			InsecureSkipTLSVerify: os.Getenv("PLUGIN_SKIP_TLS") == "false", // TODO: coerce from JSON true false into bool
		},
	}

	a := os.Getenv("PLUGIN_ADDITIONAL_ANNOTATIONS")
	if a != "" {
		var aa map[string]string
		if err := json.Unmarshal([]byte(a), &aa); err != nil {
			log.Fatalf("failed to unmarshall additional annotations: %s", err)
		}
		plugin.KubeConfig.AdditionalAnnotations = aa
	}

	fmt.Printf(os.Getenv("PLUGIN_SKIP_TLS"))
	fmt.Println("originally from danielgormly/drone-plugin-kube@0.0.2 https://github.com/danielgormly/drone-plugin-kube")
	err := plugin.Exec()
	if err != nil {
		log.Fatalf("⛔️ Fatal error: \n%s", err)
	}
}

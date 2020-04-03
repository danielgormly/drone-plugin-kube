package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	plugin := Plugin{
		Template:      os.Getenv("PLUGIN_TEMPLATE"),
		ConfigMapFile: os.Getenv("PLUGIN_CONFIGMAP_FILE"),
		KubeConfig: KubeConfig{
			Token:                 os.Getenv("PLUGIN_TOKEN"),
			Server:                os.Getenv("PLUGIN_SERVER"),
			Ca:                    os.Getenv("PLUGIN_CA"),
			Namespace:             os.Getenv("PLUGIN_NAMESPACE"),
			InsecureSkipTLSVerify: os.Getenv("PLUGIN_SKIP_TLS") == "false", // TODO: coerce from JSON true false into bool
		},
	}
	fmt.Printf(os.Getenv("PLUGIN_SKIP_TLS"))
	fmt.Println("danielgormly/drone-plugin-kube@0.2.0 https://github.com/danielgormly/drone-plugin-kube")
	err := plugin.Exec()
	if err != nil {
		log.Fatalf("⛔️ Fatal error: \n%s", err)
	}
}

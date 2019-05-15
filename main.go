package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	plugin := Plugin{
		Template: os.Getenv("PLUGIN_TEMPLATE"),
		KubeConfig: KubeConfig{
			Token:                 os.Getenv("PLUGIN_TOKEN"),
			Server:                os.Getenv("PLUGIN_SERVER"),
			Ca:                    os.Getenv("PLUGIN_CA"),
			Namespace:             os.Getenv("PLUGIN_NAMESPACE"),
			InsecureSkipTLSVerify: os.Getenv("PLUGIN_SKIP_TLS") == "false",
		},
	}
	fmt.Printf(os.Getenv("PLUGIN_SKIP_TLS"))
	err := plugin.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

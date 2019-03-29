package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Kubano"
	app.Usage = "Use with Drone CI"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ca",
			Usage:  "Certificate Authority cert to use (Base-64 encoded)",
			EnvVar: "PLUGIN_CA",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Kubernetes service token",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "server",
			Usage:  "Kubernetes server address",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "namespace to use: 'default' is the default",
			EnvVar: "PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "template file to use for deployment e.g. deployment.yaml",
			EnvVar: "PLUGIN_TEMPLATE",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}
	plugin := Plugin{
		Template: c.String("template"),
		KubeConfig: KubeConfig{
			Token:     c.String("token"),
			Server:    c.String("server"),
			Ca:        c.String("ca"),
			Namespace: c.String("namespace"),
		},
	}
	return plugin.Exec()
}

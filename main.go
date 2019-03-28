package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Kubano"
	app.Usage = "To be used within DroneCI"
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Hello %q", c.Args().Get(0))
		return nil
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ca",
			EnvVar: "APP_LANG",
			Usage:  "KUBE_CA,PLUGIN_CA",
		},
		cli.StringFlag{
			Name:   "token",
			EnvVar: "APP_LANG",
			Usage:  "KUBE_TOKEN,PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:   "server",
			EnvVar: "APP_LANG",
			Usage:  "KUBE_SERVER,PLUGIN_SERVER",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

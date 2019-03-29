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
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "namespace to use: 'default' is the default :-)",
			EnvVar: "KUBE_NAMESPACE,PLUGIN_NAMESPACE",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "template file to use for deployment e.g. deployment.yaml",
			EnvVar: "KUBE_TEMPLATE,PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	// kubernetes token

	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Tag:     c.String("build.tag"),
			Number:  c.Int("build.number"),
			Event:   c.String("build.event"),
			Status:  c.String("build.status"),
			Commit:  c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Author:  c.String("commit.author"),
			Link:    c.String("build.link"),
			Started: c.Int64("build.started"),
			Created: c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Token:     c.String("token"),
			Server:    c.String("server"),
			Ca:        c.String("ca"),
			Namespace: c.String("namespace"),
			Template:  c.String("template"),
		},
	}

	return plugin.Exec()
}

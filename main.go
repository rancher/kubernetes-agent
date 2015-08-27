package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-agent/kubernetesevents"
	"github.com/rancher/kubernetes-agent/rancherevents"
)

func main() {
	app := cli.NewApp()
	app.Name = "kubernetes-agent"
	app.Usage = "Start the Rancher kubernetes agent"
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubernetes-url",
			Value:  "http://localhost:8080",
			Usage:  "URL for kubernetes API",
			EnvVar: "KUBERNETES_URL",
		},
		cli.StringFlag{
			Name:   "cattle-url",
			Usage:  "URL for cattle API",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattle-access-key",
			Usage:  "Cattle API Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattle-secret-key",
			Usage:  "Cattle API Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:   "worker-count",
			Value:  50,
			Usage:  "Number of workers for handling events",
			EnvVar: "WORKER_COUNT",
		},
	}

	app.Run(os.Args)
}

func launch(c *cli.Context) {
	conf := config.Conf(c)

	resultChan := make(chan error)

	rClient, err := config.GetRancherClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	kClient := kubernetesclient.NewClient(conf.KubernetesURL, true)

	rcHandler := kubernetesevents.NewHandler(rClient, kClient, kubernetesevents.RCKind)
	svcHandler := kubernetesevents.NewHandler(rClient, kClient, kubernetesevents.ServiceKind)
	handlers := []kubernetesevents.Handler{rcHandler, svcHandler}

	go func(rc chan error) {
		err := kubernetesevents.ConnectToEventStream(handlers, conf)
		log.Errorf("Kubernetes stream listener exited with error: %s", err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err := rancherevents.ConnectToEventStream(conf)
		log.Errorf("Rancher stream listener exited with error: %s", err)
		rc <- err
	}(resultChan)

	<-resultChan
	log.Info("Exiting.")
}

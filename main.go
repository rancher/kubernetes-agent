package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/leodotcloud/log"
	logserver "github.com/leodotcloud/log/server"
	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/healthcheck"
	"github.com/rancher/kubernetes-agent/hostlabels"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-agent/rancherevents"
	"github.com/rancher/kubernetes-agent/watchevents"
)

func main() {
	logserver.StartServerWithDefaults()
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
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  10240,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.IntFlag{
			Name:  "host-labels-update-interval",
			Value: 5,
			Usage: "The frequency at which host labels should be updated",
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

	kClient := kubernetesclient.NewClient(conf.KubernetesURL)

	svcHandler := watchevents.NewServiceHandler(rClient, kClient)

	nsHandler := watchevents.NewNamespaceHandler(rClient, kClient)

	svcHandler.Start()
	defer svcHandler.Stop()

	nsHandler.Start()
	defer nsHandler.Stop()

	go func(rc chan error) {
		err := rancherevents.ConnectToEventStream(conf)
		log.Errorf("Rancher stream listener exited with error: %s", err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err := healthcheck.StartHealthCheck(conf.HealthCheckPort)
		log.Errorf("Rancher healthcheck exited with error: %s", err)
		rc <- err
	}(resultChan)

	go func(rc chan error) {
		err := hostlabels.StartHostLabelSync(c.Int("host-labels-update-interval"), kClient)
		log.Errorf("Rancher hostLabel sync service exited with error: %s", err)
		rc <- err
	}(resultChan)

	<-resultChan
	log.Info("Exiting.")
}

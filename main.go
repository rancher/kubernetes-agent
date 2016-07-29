package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/healthcheck"
	"github.com/rancher/kubernetes-agent/hostlabels"
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
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  10240,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.StringSliceFlag{
			Name: "watch-kind",
			Value: &cli.StringSlice{"namespaces", "services", "replicationcontrollers", "pods",
				"deployments", "ingresses", "jobs", "horizontalpodautoscalers", "persistentvolumes",
				"persistentvolumeclaims", "replicasets", "secrets"},
			Usage: "Which k8s kinds to watch and report changes to Rancher",
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

	kClient := kubernetesclient.NewClient(conf.KubernetesURL, true)

	svcHandler, err := kubernetesevents.NewHandler(rClient, kClient, kubernetesevents.ServiceKind)
	if err != nil {
		log.Fatal(err)
	}
	nsHandler, err := kubernetesevents.NewHandler(rClient, kClient, kubernetesevents.NamespaceKind)
	if err != nil {
		log.Fatal(err)
	}
	handlers := []kubernetesevents.Handler{svcHandler, nsHandler}

	log.Info("Watching changes for kinds: ", c.StringSlice("watch-kind"))
	for _, kind := range c.StringSlice("watch-kind") {
		handlers = append(handlers, kubernetesevents.NewChangeHandler(rClient, kClient, kind))
	}

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

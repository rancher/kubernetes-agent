package rancherevents

import (
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-agent/rancherevents/eventhandlers"
)

func ConnectToEventStream(rClient *client.RancherClient, conf config.Config) error {

	kClient := kubernetesclient.NewClient(conf.KubernetesURL)

	eventHandlers := map[string]revents.EventHandler{
		"compute.instance.providelabels": eventhandlers.NewProvideLablesHandler(kClient).Handler,
		"config.update":                  eventhandlers.NewPingHandler().Handler,
		"ping":                           eventhandlers.NewPingHandler().Handler,
	}

	router, err := revents.NewEventRouter(rClient, conf.WorkerCount, eventHandlers)
	if err != nil {
		return err
	}
	err = router.Start(nil)
	return err
}

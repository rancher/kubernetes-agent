package rancherevents

import (
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-agent/rancherevents/eventhandlers"
)

func ConnectToEventStream(rClient *client.RancherClient, kClient *kubernetesclient.Client, conf config.Config) error {

	eventHandlers := map[string]revents.EventHandler{
		"config.update":   eventhandlers.NewPingHandler().Handler,
		"ping":            eventhandlers.NewPingHandler().Handler,
		"host.deactivate": eventhandlers.NewHostHandler(kClient).Handler,
		"host.activate":   eventhandlers.NewHostHandler(kClient).Handler,
		"host.evacuate":   eventhandlers.NewHostHandler(kClient).Handler,
	}

	router, err := revents.NewEventRouter(rClient, conf.WorkerCount, eventHandlers)
	if err != nil {
		return err
	}
	err = router.Start(nil)
	return err
}

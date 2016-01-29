package eventhandlers

import (
	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	util "github.com/rancher/kubernetes-agent/rancherevents/util"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	if err := util.CreateAndPublishReply(event, cli); err != nil {
		return err
	}
	return nil
}

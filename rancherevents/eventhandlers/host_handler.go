package eventhandlers

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v3"
	"github.com/rancher/kubernetes-agent/hostwatch"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	util "github.com/rancher/kubernetes-agent/rancherevents/util"
)

type HostHandler struct {
	kClient *kubernetesclient.Client
}

func NewHostHandler(kClient *kubernetesclient.Client) *HostHandler {
	return &HostHandler{
		kClient: kClient,
	}
}

const (
	ActivateEvent   = "host.activate"
	DeactivateEvent = "host.deactivate"
	EvacuateEvent   = "host.evacuate"
)

var (
	hostWaitTimeout  time.Duration = 10
	hostDrainTimeout time.Duration = 60
)

func (h *HostHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	var timeout time.Duration
	if timeout = hostWaitTimeout; event.Name == EvacuateEvent {
		timeout = hostDrainTimeout
	}
	nodeMap := GetStringMap(event.Data, "host", "data", "fields")
	nodeName := nodeMap["nodeName"]
	if nodeName == "" {
		return fmt.Errorf("Failed to parse event data")
	}
	// Wait for response from k8s api
	desiredState := make(chan string, 1)
	go h.getDesiredState(nodeName, event.Name, desiredState)
	select {
	case state := <-desiredState:
		log.Infof("Updated node: [%s] with new schedulable state: [%s]", nodeName, state)
		if err := util.CreateAndPublishReply(event, cli); err != nil {
			return fmt.Errorf("Error publishing reply: %v", err)
		}
		return nil
	case <-time.After(time.Second * timeout):
		return fmt.Errorf("Timeout waiting for kubernetes node")
	}

}

func (h *HostHandler) getDesiredState(nodeName string, eventName string, desiredState chan string) {
	gotState := false
	for !gotState {
		node, err := h.kClient.Node.ByName(nodeName)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", nodeName, err)
			return
		}
		if eventName == ActivateEvent {
			if !node.Spec.Unschedulable {
				desiredState <- "Schedulable"
				gotState = true
			}
		} else if eventName == DeactivateEvent {
			if node.Spec.Unschedulable {
				desiredState <- "Unschedulable"
				gotState = true
			}
		} else if eventName == EvacuateEvent {
			if state, ok := node.ObjectMeta.Labels[hostwatch.DrainLabelName]; ok && (state == hostwatch.DrainLabelValue) {
				desiredState <- "Evacuated"
				gotState = true
			}
		}
	}
}

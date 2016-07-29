package kubernetesevents

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

const NamespaceKind string = "namespaces"
const namespaceEventTypePrefix string = "stack."

// Capable of handling RC and Service events
type NamespaceHandler struct {
	rancherClient *client.RancherClient
	kClient       *kubernetesclient.Client
	kindHandled   string
}

func (h *NamespaceHandler) GetKindHandled() string {
	return h.kindHandled
}

func (h *NamespaceHandler) Handle(event model.WatchEvent) error {
	if i, ok := event.Object.(map[string]interface{}); ok {
		var ns model.Namespace
		mapstructure.Decode(i, &ns)
		if ns == (model.Namespace{}) || ns.Spec == nil {
			log.Errorf("Could not decode %+v to namespace", i)
			return nil
		}
		Uid := ns.Metadata.Uid

		var serviceEvent *client.ExternalServiceEvent

		switch event.Type {
		case "DELETED":
			serviceEvent = &client.ExternalServiceEvent{
				ExternalId: "kubernetes://" + Uid,
				EventType:  namespaceEventTypePrefix + "remove",
			}
			serviceEvent.Environment = &client.Environment{
				Kind: "environment",
			}
			serviceEvent.Service = client.Service{
				Kind: "kubernetesService",
			}
		default:
			return nil
		}
		_, err := h.rancherClient.ExternalServiceEvent.Create(serviceEvent)
		return err
	}
	return fmt.Errorf("Could not decode event %+v", event)
}

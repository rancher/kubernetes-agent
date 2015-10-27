package kubernetesevents

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

const RCKind string = "replicationcontrollers"
const ServiceKind string = "services"

func NewHandler(rancherClient *client.RancherClient, kubernetesClient *kubernetesclient.Client, kindHandled string) *GenericHandler {
	return &GenericHandler{
		rancherClient: rancherClient,
		kClient:       kubernetesClient,
		kindHandled:   kindHandled,
	}
}

// Capable of handling RC and Service events
type GenericHandler struct {
	rancherClient *client.RancherClient
	kClient       *kubernetesclient.Client
	kindHandled   string
}

func (h *GenericHandler) GetKindHandled() string {
	return h.kindHandled
}

func (h *GenericHandler) Handle(event model.WatchEvent) error {

	if i, ok := event.Object.(map[string]interface{}); ok {
		var metadata *model.ObjectMeta
		var kind string
		var selector map[string]interface{}
		if h.kindHandled == RCKind {
			var rc model.ReplicationController
			mapstructure.Decode(i, &rc)
			kind = rc.Kind
			selector = rc.Spec.Selector
			metadata = rc.Metadata
		} else if h.kindHandled == ServiceKind {
			var svc model.Service
			mapstructure.Decode(i, &svc)
			kind = svc.Kind
			selector = svc.Spec.Selector
			metadata = svc.Metadata
		} else {
			return fmt.Errorf("Unrecognized handled kind [%s].", h.kindHandled)
		}

		ts, err := time.Parse(time.RFC3339, metadata.CreationTimestamp)
		if err != nil {
			return err
		}

		serviceEvent := &client.ServiceEvent{
			Name:              metadata.Name,
			ExternalId:        metadata.Uid,
			EventType:         constructEventType(event),
			ResourceType:      constructResourceType(kind),
			ExternalTimestamp: ts.Unix(),
		}

		switch event.Type {
		case "MODIFIED":
			fallthrough

		case "ADDED":
			err = h.add(selector, metadata, event, serviceEvent, serviceEvent.ResourceType)
			if err != nil {
				return err
			}

		case "DELETED":
			// No extra fields required for delete
		default:
			return nil
		}

		_, err = h.rancherClient.ServiceEvent.Create(serviceEvent)
		return err

	}
	return fmt.Errorf("Couldn't decode event [%#v]", event)
}

func (h *GenericHandler) add(selectorMap map[string]interface{}, metadata *model.ObjectMeta, event model.WatchEvent, serviceEvent *client.ServiceEvent, kind string) error {
	cattleServiceFields := make(map[string]interface{})
	serviceEvent.CattleResourceFields = cattleServiceFields

	cattleServiceFields["template"] = event.Object
	cattleServiceFields["kind"] = kind
	cattleServiceFields["name"] = metadata.Name
	cattleServiceFields["externalId"] = metadata.Uid

	var buffer bytes.Buffer
	for key, v := range selectorMap {
		if val, ok := v.(string); ok {
			buffer.WriteString(key)
			buffer.WriteString("=")
			buffer.WriteString(val)
			buffer.WriteString(",")
		}
	}
	selector := buffer.String()
	selector = strings.TrimSuffix(selector, ",")
	cattleServiceFields["selectorContainer"] = selector

	namespace, err := h.kClient.Namespace.ByName(metadata.Namespace)
	if err != nil {
		return err
	}
	env := make(map[string]string)
	env["name"] = namespace.Metadata.Name
	env["externalId"] = namespace.Metadata.Uid
	serviceEvent.Environment = env

	return err
}

func constructEventType(event model.WatchEvent) string {
	return strings.ToLower("externalprovider.service." + event.Type)
}

func constructResourceType(kind string) string {
	return strings.ToLower("kubernetes" + kind)
}

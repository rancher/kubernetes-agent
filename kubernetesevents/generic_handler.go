package kubernetesevents

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

const RCKind string = "replicationcontrollers"
const ServiceKind string = "services"
const eventTypePrefix string = "service."

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
			if rc == (model.ReplicationController{}) || rc.Spec == nil {
				log.Infof("Couldn't decode %+v to rc.", i)
				return nil
			}
			kind = rc.Kind
			selector = rc.Spec.Selector
			metadata = rc.Metadata
		} else if h.kindHandled == ServiceKind {
			var svc model.Service
			mapstructure.Decode(i, &svc)
			if svc == (model.Service{}) || svc.Spec == nil {
				log.Infof("Couldn't decode %+v to service.", i)
				return nil
			}
			kind = svc.Kind
			selector = svc.Spec.Selector
			metadata = svc.Metadata
		} else {
			return fmt.Errorf("Unrecognized handled kind [%s].", h.kindHandled)
		}

		serviceEvent := &client.ExternalServiceEvent{
			ExternalId: metadata.Uid,
			EventType:  constructEventType(event),
		}

		switch event.Type {
		case "MODIFIED":
			fallthrough

		case "ADDED":
			err := h.add(selector, metadata, event, serviceEvent, constructResourceType(kind))
			if err != nil {
				return err
			}

		case "DELETED":
			service := client.Service{
				Kind: constructResourceType(kind),
			}
			serviceEvent.Service = service
		default:
			return nil
		}

		_, err := h.rancherClient.ExternalServiceEvent.Create(serviceEvent)
		return err

	}
	return fmt.Errorf("Couldn't decode event [%#v]", event)
}

func (h *GenericHandler) add(selectorMap map[string]interface{}, metadata *model.ObjectMeta, event model.WatchEvent, serviceEvent *client.ExternalServiceEvent, kind string) error {
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

	fields := map[string]interface{}{"template": event.Object}
	data := map[string]interface{}{"fields": fields}

	service := client.Service{
		Kind:              kind,
		Name:              metadata.Name,
		ExternalId:        metadata.Uid,
		SelectorContainer: selector,
		Data:              data,
	}
	serviceEvent.Service = service

	env := make(map[string]string)

	if metadata.Namespace == "kube-system" {
		env["name"] = metadata.Namespace
		env["externalId"] = "kubernetes://" + metadata.Namespace
	} else {
		namespace, err := h.kClient.Namespace.ByName(metadata.Namespace)
		if err != nil {
			return err
		}
		env["name"] = namespace.Metadata.Name
		env["externalId"] = "kubernetes://" + namespace.Metadata.Uid
	}
	serviceEvent.Environment = env

	return nil
}

func constructEventType(event model.WatchEvent) string {
	switch strings.ToLower(event.Type) {
	case "added":
		return eventTypePrefix + "create"
	case "modified":
		return eventTypePrefix + "update"
	case "deleted":
		return eventTypePrefix + "remove"
	default:
		return eventTypePrefix + event.Type
	}
}

func constructResourceType(kind string) string {
	return "kubernetes" + kind
}

package kubernetesevents

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"

	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

const (
	kubernetesServiceKind = "kubernetesService"
)

type SyncHandler interface {
	Add(interface{}) error
	Delete(interface{}) error
	Decode(model.WatchEvent) (interface{}, error)
	GetListURL() string
	GetWatchURL() string
	GetKey(model.WatchEvent) (string, error)
}

type serviceHandler struct {
	rClient *client.RancherClient
	kClient *kubernetesclient.Client
	baseURL string
}

func NewServiceHandler(rClient *client.RancherClient, kClient *kubernetesclient.Client, conf config.Config) *serviceHandler {
	sHandler := &serviceHandler{
		rClient: rClient,
		kClient: kClient,
		baseURL: conf.KubernetesURL,
	}
	return sHandler
}

func (s *serviceHandler) Decode(event model.WatchEvent) (interface{}, error) {
	i, ok := event.Object.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("Couldn't decode service event [#v]", event)
	}

	var svc model.Service
	mapstructure.Decode(i, &svc)
	if svc == (model.Service{}) || svc.Spec == nil {
		log.Infof("Couldn't decode %+v to service.", i)
		return nil, fmt.Errorf("Service object is empty")
	}
	return svc, nil
}

func (s *serviceHandler) GetKey(event model.WatchEvent) (string, error) {
	val, err := s.Decode(event)
	if err != nil {
		log.Errorf("Error computing key for event %v", err)
		return "", err
	}
	return val.(model.Service).Metadata.Uid, nil
}

func (s *serviceHandler) Add(svc interface{}) error {
	realSVC := svc.(model.Service)

	kind := kubernetesServiceKind
	metadata := realSVC.Metadata
	selectorMap := realSVC.Spec.Selector
	clusterIp := realSVC.Spec.ClusterIP

	var serviceEvent = &client.ExternalServiceEvent{}
	serviceEvent.ExternalId = metadata.Uid
	serviceEvent.EventType = "service.create"

	if selectorMap != nil {
		selectorMap["io.kubernetes.pod.namespace"] = metadata.Namespace
	}

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

	fields := map[string]interface{}{"template": realSVC}
	data := map[string]interface{}{"fields": fields}

	rancherUuid, _ := metadata.Labels["io.rancher.uuid"].(string)
	var vip string
	if !strings.EqualFold(clusterIp, "None") {
		vip = clusterIp
	}
	service := client.Service{
		Kind:              kind,
		Name:              metadata.Name,
		ExternalId:        metadata.Uid,
		SelectorContainer: selector,
		Data:              data,
		Uuid:              rancherUuid,
		Vip:               vip,
	}
	serviceEvent.Service = service

	env := make(map[string]string)

	if metadata.Namespace == "kube-system" {
		env["name"] = metadata.Namespace
		env["externalId"] = "kubernetes://" + metadata.Namespace
	} else {
		namespace, err := s.kClient.Namespace.ByName(metadata.Namespace)
		if err != nil {
			return err
		}
		env["name"] = namespace.Metadata.Name
		env["externalId"] = "kubernetes://" + namespace.Metadata.Uid
		rancherUuid, _ := namespace.Metadata.Labels["io.rancher.uuid"].(string)
		env["uuid"] = rancherUuid
	}
	serviceEvent.Environment = env
	_, err := s.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (s *serviceHandler) Delete(svc interface{}) error {
	realSVC := svc.(model.Service)

	kind := kubernetesServiceKind
	metadata := realSVC.Metadata

	var serviceEvent = &client.ExternalServiceEvent{}
	serviceEvent.ExternalId = metadata.Uid
	serviceEvent.EventType = "service.remove"
	service := client.Service{
		Kind: kind,
	}
	serviceEvent.Service = service

	_, err := s.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (s *serviceHandler) GetListURL() string {
	return fmt.Sprintf("%s/api/v1/services", s.baseURL)
}

func (s *serviceHandler) GetWatchURL() string {
	baseURL := strings.Replace(s.GetListURL(), "http", "ws", 1)
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Couldn't parse URL, err: %v", err)
	}
	q := u.Query()
	q.Set("watch", "true")
	u.RawQuery = q.Encode()
	return u.String()
}

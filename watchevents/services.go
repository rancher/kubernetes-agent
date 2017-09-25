package watchevents

import (
	"bytes"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

const (
	kubernetesServiceKind = "kubernetesService"
)

type serviceHandler struct {
	rClient          *client.RancherClient
	kClient          *kubernetesclient.Client
	serviceWatchChan chan struct{}
}

func NewServiceHandler(rClient *client.RancherClient, kClient *kubernetesclient.Client) *serviceHandler {
	sHandler := &serviceHandler{
		rClient: rClient,
		kClient: kClient,
	}
	return sHandler
}

func (s *serviceHandler) startServiceWatch() chan struct{} {
	watchlist := cache.NewListWatchFromClient(s.kClient.K8sClientSet.Core().RESTClient(), "services", v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Received event: [ADDED] for Service: %s, Handling Add event.", key)
				err := s.add(obj, "Added")
				if err != nil {
					logrus.Errorf("Error Handling event: [ADDED] for Service: %v", err)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Received event: [DELETED] for Service: %s, Handling Delete event.", key)
				err := s.delete(obj)
				if err != nil {
					logrus.Errorf("Error Handling event: [DELETED] for Service: %v", err)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(newObj)
				logrus.Infof("Received event: [MODIFIED] for Service: %s, Handling Modified event.", key)
				err := s.add(newObj, "Modified")
				if err != nil {
					logrus.Errorf("Error Handling event: [MODIFIED] for Service: %v", err)
				}
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	return stop
}

func (s *serviceHandler) add(svc interface{}, eventType string) error {
	realSVC := svc.(*v1.Service)
	kind := kubernetesServiceKind
	metadata := realSVC.ObjectMeta
	selectorMap := realSVC.Spec.Selector
	clusterIP := realSVC.Spec.ClusterIP

	var serviceEvent = &client.ExternalServiceEvent{}
	serviceEvent.ExternalId = string(metadata.UID)
	if eventType == "Added" {
		serviceEvent.EventType = "service.create"
	} else {
		serviceEvent.EventType = "service.update"
	}

	if selectorMap != nil {
		selectorMap["io.kubernetes.pod.namespace"] = metadata.Namespace
	}

	var buffer bytes.Buffer
	for key, v := range selectorMap {
		buffer.WriteString(key)
		buffer.WriteString("=")
		buffer.WriteString(v)
		buffer.WriteString(",")
	}
	selector := buffer.String()
	selector = strings.TrimSuffix(selector, ",")

	fields := map[string]interface{}{"template": realSVC}
	data := map[string]interface{}{"fields": fields}

	rancherUUID, _ := metadata.Labels["io.rancher.uuid"]
	var vip string
	if !strings.EqualFold(clusterIP, "None") {
		vip = clusterIP
	}
	service := client.Service{
		Kind:              kind,
		Name:              metadata.Name,
		ExternalId:        string(metadata.UID),
		SelectorContainer: selector,
		Data:              data,
		Uuid:              rancherUUID,
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
		env["name"] = namespace.Name
		env["externalId"] = "kubernetes://" + string(namespace.UID)
		rancherUUID, _ := namespace.Labels["io.rancher.uuid"]
		env["uuid"] = rancherUUID
	}
	serviceEvent.Environment = env
	_, err := s.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (s *serviceHandler) delete(svc interface{}) error {
	realSVC := svc.(*v1.Service)
	kind := kubernetesServiceKind
	metadata := realSVC.ObjectMeta

	var serviceEvent = &client.ExternalServiceEvent{}
	serviceEvent.ExternalId = string(metadata.UID)
	serviceEvent.EventType = "service.remove"
	service := client.Service{
		Kind: kind,
	}
	serviceEvent.Service = service

	_, err := s.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (s *serviceHandler) Start() {
	logrus.Infof("Starting service watch")
	s.serviceWatchChan = s.startServiceWatch()
}

func (s *serviceHandler) Stop() {
	if s.serviceWatchChan != nil {
		s.serviceWatchChan <- struct{}{}
	}
}

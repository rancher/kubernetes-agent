package watchevents

import (
	"bytes"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v3"
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
	watchlist := cache.NewListWatchFromClient(s.kClient.K8sClient.Core().RESTClient(), "services", v1.NamespaceAll, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Service{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: serviceAddDelete(func(service v1.Service) {
				logrus.Infof("Received event: [%s] for Service: %s/%s, Handling event.", AddedEventType, service.Namespace, service.Name)
				if err := s.add(service, AddedEventType); err != nil {
					logrus.Errorf("Error Handling event: [%s] for Service %s/%s: %v", AddedEventType, service.Namespace, service.Name, err)
				}
			}),
			DeleteFunc: serviceAddDelete(func(service v1.Service) {
				logrus.Infof("Received event: [%s] for Service: %s/%s, Handling event.", DeletedEventType, service.Namespace, service.Name)
				if err := s.delete(service); err != nil {
					logrus.Errorf("Error Handling event: [%s] for Service %s/%s: %v", DeletedEventType, service.Namespace, service.Name, err)
				}
			}),
			UpdateFunc: serviceModify(func(service v1.Service) {
				logrus.Infof("Received event: [%s] for Service: %s/%s, Handling event.", ModifiedEventType, service.Namespace, service.Name)
				if err := s.add(service, ModifiedEventType); err != nil {
					logrus.Errorf("Error Handling event: [%s] for Service %s/%s: %v", ModifiedEventType, service.Namespace, service.Name, err)
				}
			}),
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	return stop
}

func (s *serviceHandler) add(realSVC v1.Service, eventType string) error {
	kind := kubernetesServiceKind
	metadata := realSVC.ObjectMeta
	selectorMap := realSVC.Spec.Selector
	clusterIP := realSVC.Spec.ClusterIP

	var serviceEvent = &client.ExternalServiceEvent{}
	serviceEvent.ExternalId = string(metadata.UID)
	if eventType == AddedEventType {
		serviceEvent.EventType = "service.create"
	} else if eventType == ModifiedEventType {
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

	rancherUUID, _ := metadata.Labels["io.rancher.uuid"]
	var vip string
	if !strings.EqualFold(clusterIP, "None") {
		vip = clusterIP
	}
	service := client.Service{
		Kind:       kind,
		Name:       metadata.Name,
		ExternalId: string(metadata.UID),
		Selector:   selector,
		Uuid:       rancherUUID,
		Vip:        vip,
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

func (s *serviceHandler) delete(realSVC v1.Service) error {
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

func serviceAddDelete(f func(v1.Service)) func(interface{}) {
	return func(obj interface{}) {
		service := obj.(*v1.Service)
		f(*service)
	}
}

func serviceModify(f func(v1.Service)) func(interface{}, interface{}) {
	return func(oldObj, newObj interface{}) {
		serviceAddDelete(f)(newObj)
	}
}

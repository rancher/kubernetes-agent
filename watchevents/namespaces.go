package watchevents

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

type namespaceHandler struct {
	rClient     *client.RancherClient
	kClient     *kubernetesclient.Client
	nsWatchChan chan struct{}
}

func NewNamespaceHandler(rClient *client.RancherClient, kClient *kubernetesclient.Client) *namespaceHandler {
	nsHandler := &namespaceHandler{
		rClient: rClient,
		kClient: kClient,
	}
	return nsHandler
}

func (n *namespaceHandler) startNamespaceWatch() chan struct{} {
	watchlist := cache.NewListWatchFromClient(n.kClient.K8sClientSet.Core().RESTClient(), "namespaces", "", fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Namespace{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Skipping event: [ADDED] for namespace: %s", key)
			},
			DeleteFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Received event: [DELETED] for Namespace: %s, Handling Delete event.", key)
				err := n.delete(obj, "Deleted")
				if err != nil {
					logrus.Errorf("Error Handling event: [DELETED] for namespace: %v", err)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(newObj)
				logrus.Infof("Skipping event: [MODIFIED] for namespace: %s", key)
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	return stop
}

func (n *namespaceHandler) delete(ns interface{}, eventType string) error {
	realNS := ns.(*v1.Namespace)
	var metadata metav1.ObjectMeta
	var kind string
	var prefix string
	var serviceEvent = &client.ExternalServiceEvent{}
	prefix = "kubernetes://"
	metadata = realNS.ObjectMeta
	kind = "kubernetesService"
	serviceEvent.Environment = &client.Stack{
		Kind: "environment",
	}
	serviceEvent.ExternalId = prefix + string(metadata.UID)
	if eventType == "Deleted" {
		serviceEvent.EventType = "stack.remove"
	}
	service := client.Service{
		Kind: kind,
	}
	serviceEvent.Service = service

	_, err := n.rClient.ExternalServiceEvent.Create(serviceEvent)
	return err
}

func (n *namespaceHandler) Start() {
	logrus.Infof("Starting namespace watch")
	n.nsWatchChan = n.startNamespaceWatch()
}

func (n *namespaceHandler) Stop() {
	if n.nsWatchChan != nil {
		n.nsWatchChan <- struct{}{}
	}
}

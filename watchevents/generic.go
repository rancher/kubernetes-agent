package watchevents

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/pkg/apis/autoscaling"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	extv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type genericHandler struct {
	rClient          *client.RancherClient
	kClient          *kubernetesclient.Client
	genericWatchChan []chan struct{}
}

func NewGenericHandler(rClient *client.RancherClient, kClient *kubernetesclient.Client) *genericHandler {
	genHandler := &genericHandler{
		rClient: rClient,
		kClient: kClient,
	}
	return genHandler
}

func (g *genericHandler) startGenericWatch(resourceName string) chan struct{} {
	runtimeObj, restClient := g.getRuntimeObject(resourceName)
	watchlist := cache.NewListWatchFromClient(restClient, resourceName, "", fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		runtimeObj,
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Received event: [ADDED] for %s: %v, publishing event to Rancher", resourceName, key)
				g.rancherPatch(obj, "ADDED")
			},
			DeleteFunc: func(obj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(obj)
				logrus.Infof("Received event: [DELETED] for %s: %v, publishing event to Rancher", resourceName, key)
				g.rancherPatch(obj, "DELETED")
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				key, _ := cache.MetaNamespaceKeyFunc(newObj)
				logrus.Infof("Received event: [MODIFIED] for %s: %v, publishing event to Rancher", resourceName, key)
				g.rancherPatch(newObj, "MODIFIED")
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	return stop
}

func (g *genericHandler) Start(resNames []string) {
	logrus.Infof("Starting Generic watch")
	g.genericWatchChan = make([]chan struct{}, len(resNames))
	for k, v := range resNames {
		g.genericWatchChan[k] = g.startGenericWatch(v)
	}
}

func (g *genericHandler) Stop() {
	for _, channel := range g.genericWatchChan {
		if channel != nil {
			channel <- struct{}{}
		}
	}
}

func (g *genericHandler) rancherPatch(res interface{}, eventType string) error {
	_, err := g.rClient.Publish.Create(&client.Publish{
		Name: "service.kubernetes.change",
		Data: map[string]interface{}{
			"type":   eventType,
			"object": res,
		},
	})

	return err
}

func (g *genericHandler) getRuntimeObject(resName string) (runtime.Object, rest.Interface) {
	switch resName {
	case "pods":
		return &v1.Pod{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "namespaces":
		return &v1.Namespace{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "secrets":
		return &v1.Secret{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "replicationcontrollers":
		return &v1.ReplicationController{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "services":
		return &v1.Service{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "persistentvolumes":
		return &v1.PersistentVolume{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "persistentvolumeclaims":
		return &v1.PersistentVolumeClaim{}, g.kClient.K8sClientSet.Core().RESTClient()
	case "deployments":
		return &v1beta1.Deployment{}, g.kClient.K8sClientSet.AppsV1beta1().RESTClient()
	// replicasets restclient not available in 4.0.0
	case "replicasets":
		return &extv1beta1.ReplicaSet{}, g.kClient.K8sClientSet.ExtensionsV1beta1().RESTClient()
	case "jobs":
		return &batchv1.Job{}, g.kClient.K8sClientSet.BatchV1().RESTClient()
	case "ingresses":
		return &extv1beta1.Ingress{}, g.kClient.K8sClientSet.ExtensionsV1beta1().RESTClient()
	case "horizontalpodautoscalers":
		return &autoscaling.HorizontalPodAutoscaler{}, g.kClient.K8sClientSet.AutoscalingV1().RESTClient()
	default:
		return &v1.Pod{}, g.kClient.K8sClientSet.Core().RESTClient()
	}
}

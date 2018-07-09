package kubernetesclient

import (
	"github.com/rancher/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	caLocation         = "/etc/kubernetes/ssl/ca.pem"
	kubeconfigLocation = "/etc/kubernetes/ssl/kubeconfig"
)

func GetK8sClientSet(apiURL string) *kubernetes.Clientset {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigLocation)
	if apiURL != "" {
		config.Host = apiURL
	}
	if err != nil {
		log.Fatalf("Can't Build kubernetes config: %v", err)
	}
	K8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Can't initiate kubernetes client: %v", err)
	}
	return K8sClientSet
}

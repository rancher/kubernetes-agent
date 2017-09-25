package kubernetesclient

import (
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
		panic(err.Error())
	}
	K8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return K8sClientSet
}

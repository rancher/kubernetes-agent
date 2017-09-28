package kubernetesclient

import (
	"os"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeconfigLocation = "/etc/kubernetes/ssl/kubeconfig"
	inClusterConfig    = "INCLUSTER_CONFIG"
)

func GetK8sClientSet(apiURL string) *kubernetes.Clientset {
	var config *rest.Config
	var err error

	// used for backward compatibility for unit tests
	if inClusterEnv := os.Getenv(inClusterConfig); inClusterEnv == "false" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigLocation)
	} else {
		config, err = rest.InClusterConfig()
	}
	if apiURL != "" {
		config.Host = apiURL
	}
	if err != nil {
		logrus.Fatalf("Can't Build kubernetes config: %v", err)
	}
	K8sClientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Can't initiate kubernetes client: %v", err)
	}
	return K8sClientSet
}

package kubernetesclient

import "k8s.io/client-go/kubernetes"

func NewClient(apiURL string, debug bool) *Client {
	client := &Client{
		baseClient: baseClient{
			K8sClientSet: GetK8sClientSet(apiURL),
			debug:        debug,
		},
	}

	client.Pod = newPodClient(client)
	client.Namespace = newNamespaceClient(client)
	client.Service = newServiceClient(client)
	client.Node = newNodeClient(client)

	return client
}

type Client struct {
	baseClient
	Pod       PodOperations
	Namespace NamespaceOperations
	Service   ServiceOperations
	Node      NodeOperations
}

type baseClient struct {
	K8sClientSet *kubernetes.Clientset
	debug        bool
}

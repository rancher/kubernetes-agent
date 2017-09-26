package kubernetesclient

import "k8s.io/client-go/kubernetes"

func NewClient(apiURL string) *Client {
	client := &Client{
		K8sClient: GetK8sClientSet(apiURL),
	}

	client.Pod = newPodClient(client)
	client.Namespace = newNamespaceClient(client)
	client.Service = newServiceClient(client)
	client.Node = newNodeClient(client)

	return client
}

type Client struct {
	K8sClient *kubernetes.Clientset
	Pod       PodOperations
	Namespace NamespaceOperations
	Service   ServiceOperations
	Node      NodeOperations
}

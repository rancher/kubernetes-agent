package kubernetesclient

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type NodeOperations interface {
	ByName(name string) (*v1.Node, error)
	ReplaceNode(resource *v1.Node) (*v1.Node, error)
	GetNodeList() (*v1.NodeList, error)
}

func newNodeClient(client *Client) *NodeClient {
	return &NodeClient{
		client: client,
	}
}

type NodeClient struct {
	client *Client
}

func (c *NodeClient) ByName(name string) (*v1.Node, error) {
	return c.client.K8sClient.CoreV1().Nodes().Get(name, metav1.GetOptions{})
}

func (c *NodeClient) ReplaceNode(resource *v1.Node) (*v1.Node, error) {
	return c.client.K8sClient.CoreV1().Nodes().Update(resource)
}

func (c *NodeClient) GetNodeList() (*v1.NodeList, error) {
	return c.client.K8sClient.CoreV1().Nodes().List(metav1.ListOptions{})
}

package kubernetesclient

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type NodeOperations interface {
	ByName(name string) (*v1.Node, error)
	ReplaceNode(resource *v1.Node) (*v1.Node, error)
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
	resp, err := c.client.K8sClientSet.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	return resp, err
}

func (c *NodeClient) ReplaceNode(resource *v1.Node) (*v1.Node, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Nodes().Update(resource)
	return resp, err
}

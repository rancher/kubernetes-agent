package kubernetesclient

import (
	"fmt"

	"github.com/rancher/kubernetes-model/model"
)

const NodePath string = "/api/v1/nodes"
const NodeByNamePath string = "/api/v1/nodes/%s"

type NodeOperations interface {
	ByName(name string) (*model.Node, error)
	CreateNode(resource *model.Node) (*model.Node, error)
	ReplaceNode(resource *model.Node) (*model.Node, error)
	DeleteNode(name string) (*model.Status, error)
}

func newNodeClient(client *Client) *NodeClient {
	return &NodeClient{
		client: client,
	}
}

type NodeClient struct {
	client *Client
}

func (c *NodeClient) ByName(name string) (*model.Node, error) {
	resp := &model.Node{}
	path := fmt.Sprintf(NodeByNamePath, name)
	err := c.client.doGet(path, resp)
	return resp, err
}

func (c *NodeClient) CreateNode(resource *model.Node) (*model.Node, error) {
	resp := &model.Node{}
	path := fmt.Sprintf(NodePath)
	err := c.client.doPost(path, resource, resp)
	return resp, err
}

func (c *NodeClient) ReplaceNode(resource *model.Node) (*model.Node, error) {
	resp := &model.Node{}
	path := fmt.Sprintf(NodeByNamePath, resource.Metadata.Name)
	err := c.client.doPut(path, resource, resp)
	return resp, err
}

func (c *NodeClient) DeleteNode(name string) (*model.Status, error) {
	status := &model.Status{}
	path := fmt.Sprintf(NodeByNamePath, name)
	err := c.client.doDelete(path, status)
	return status, err
}

package kubernetesclient

import (
	"fmt"

	"github.com/rancher/kubernetes-model/model"
)

const NamespacePath string = "/api/v1/namespaces/"
const NamespaceByNamePath string = "/api/v1/namespaces/%s"

type NamespaceOperations interface {
	ByName(name string) (*model.Namespace, error)
	CreateNamespace(resource *model.Namespace) (*model.Namespace, error)
	ReplaceNamespace(namespace string, resource *model.Namespace) (*model.Namespace, error)
	DeleteNamespace(namespace string) (*model.Status, error)
}

func newNamespaceClient(client *Client) *NamespaceClient {
	return &NamespaceClient{
		client: client,
	}
}

type NamespaceClient struct {
	client *Client
}

func (c *NamespaceClient) ByName(name string) (*model.Namespace, error) {
	resp := &model.Namespace{}
	path := fmt.Sprintf(NamespaceByNamePath, name)
	err := c.client.doGet(path, resp)
	return resp, err
}

func (c *NamespaceClient) CreateNamespace(resource *model.Namespace) (*model.Namespace, error) {
	resp := &model.Namespace{}
	err := c.client.doPost(NamespacePath, resource, resp)
	return resp, err
}

func (c *NamespaceClient) ReplaceNamespace(name string, resource *model.Namespace) (*model.Namespace, error) {
	resp := &model.Namespace{}
	path := fmt.Sprintf(NamespaceByNamePath, name)
	err := c.client.doPut(path, resource, resp)
	return resp, err
}

func (c *NamespaceClient) DeleteNamespace(name string) (*model.Status, error) {
	status := &model.Status{}
	path := fmt.Sprintf(NamespaceByNamePath, name)
	err := c.client.doDelete(path, status)
	return status, err
}

package kubernetesclient

import (
	"fmt"

	"github.com/rancher/kubernetes-model/model"
)

const ServicePath string = "/api/v1/namespaces/%s/services"
const ServiceByNamePath string = "/api/v1/namespaces/%s/services/%s"

type ServiceOperations interface {
	ByName(namespace string, name string) (*model.Service, error)
	CreateService(namespace string, resource *model.Service) (*model.Service, error)
	ReplaceService(namespace string, resource *model.Service) (*model.Service, error)
	DeleteService(namespace string, name string) (*model.Status, error)
}

func newServiceClient(client *Client) *ServiceClient {
	return &ServiceClient{
		client: client,
	}
}

type ServiceClient struct {
	client *Client
}

func (c *ServiceClient) ByName(namespace string, name string) (*model.Service, error) {
	resp := &model.Service{}
	path := fmt.Sprintf(ServiceByNamePath, namespace, name)
	err := c.client.doGet(path, resp)
	return resp, err
}

func (c *ServiceClient) CreateService(namespace string, resource *model.Service) (*model.Service, error) {
	resp := &model.Service{}
	path := fmt.Sprintf(ServicePath, namespace)
	err := c.client.doPost(path, resource, resp)
	return resp, err
}

func (c *ServiceClient) ReplaceService(namespace string, resource *model.Service) (*model.Service, error) {
	resp := &model.Service{}
	path := fmt.Sprintf(ServiceByNamePath, namespace, resource.Metadata.Name)
	err := c.client.doPut(path, resource, resp)
	return resp, err
}

func (c *ServiceClient) DeleteService(namespace string, name string) (*model.Status, error) {
	status := &model.Status{}
	path := fmt.Sprintf(ServiceByNamePath, namespace, name)
	err := c.client.doDelete(path, status)
	return status, err
}

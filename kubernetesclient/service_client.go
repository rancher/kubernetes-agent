package kubernetesclient

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type ServiceOperations interface {
	ByName(namespace string, name string) (*v1.Service, error)
	CreateService(namespace string, resource *v1.Service) (*v1.Service, error)
	ReplaceService(namespace string, resource *v1.Service) (*v1.Service, error)
	DeleteService(namespace string, name string) error
}

func newServiceClient(client *Client) *ServiceClient {
	return &ServiceClient{
		client: client,
	}
}

type ServiceClient struct {
	client *Client
}

func (c *ServiceClient) ByName(namespace string, name string) (*v1.Service, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	return resp, err
}

func (c *ServiceClient) CreateService(namespace string, resource *v1.Service) (*v1.Service, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Services(namespace).Create(resource)
	return resp, err
}

func (c *ServiceClient) ReplaceService(namespace string, resource *v1.Service) (*v1.Service, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Services(namespace).Update(resource)
	return resp, err
}

func (c *ServiceClient) DeleteService(namespace string, name string) error {
	err := c.client.K8sClientSet.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
	return err
}

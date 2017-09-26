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
	return c.client.K8sClient.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func (c *ServiceClient) CreateService(namespace string, resource *v1.Service) (*v1.Service, error) {
	return c.client.K8sClient.CoreV1().Services(namespace).Create(resource)
}

func (c *ServiceClient) ReplaceService(namespace string, resource *v1.Service) (*v1.Service, error) {
	return c.client.K8sClient.CoreV1().Services(namespace).Update(resource)
}

func (c *ServiceClient) DeleteService(namespace string, name string) error {
	return c.client.K8sClient.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{})
}

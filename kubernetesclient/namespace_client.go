package kubernetesclient

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type NamespaceOperations interface {
	ByName(name string) (*v1.Namespace, error)
	CreateNamespace(resource *v1.Namespace) (*v1.Namespace, error)
	DeleteNamespace(namespace string) error
}

func newNamespaceClient(client *Client) *NamespaceClient {
	return &NamespaceClient{
		client: client,
	}
}

type NamespaceClient struct {
	client *Client
}

func (c *NamespaceClient) ByName(name string) (*v1.Namespace, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	return resp, err
}

func (c *NamespaceClient) CreateNamespace(resource *v1.Namespace) (*v1.Namespace, error) {
	resp, err := c.client.K8sClientSet.CoreV1().Namespaces().Create(resource)
	return resp, err
}

func (c *NamespaceClient) DeleteNamespace(name string) error {
	err := c.client.K8sClientSet.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	return err
}

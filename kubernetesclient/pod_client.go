package kubernetesclient

import (
	"fmt"

	"github.com/rancher/kubernetes-model/model"
)

const PodPath string = "/api/v1/namespaces/%s/pods"
const PodByNamePath string = "/api/v1/namespaces/%s/pods/%s"

type PodOperations interface {
	ByName(namespace string, name string) (*model.Pod, error)
	CreatePod(namespace string, resource *model.Pod) (*model.Pod, error)
	ReplacePod(namespace string, resource *model.Pod) (*model.Pod, error)
	DeletePod(namespace string, name string) (*model.Status, error)
}

func newPodClient(client *Client) *PodClient {
	return &PodClient{
		client: client,
	}
}

type PodClient struct {
	client *Client
}

func (c *PodClient) ByName(namespace string, name string) (*model.Pod, error) {
	resp := &model.Pod{}
	path := fmt.Sprintf(PodByNamePath, namespace, name)
	err := c.client.doGet(path, resp)
	return resp, err
}

func (c *PodClient) CreatePod(namespace string, resource *model.Pod) (*model.Pod, error) {
	resp := &model.Pod{}
	path := fmt.Sprintf(PodPath, namespace)
	err := c.client.doPost(path, resource, resp)
	return resp, err
}

func (c *PodClient) ReplacePod(namespace string, resource *model.Pod) (*model.Pod, error) {
	resp := &model.Pod{}
	path := fmt.Sprintf(PodByNamePath, namespace, resource.Metadata.Name)
	err := c.client.doPut(path, resource, resp)
	return resp, err
}

func (c *PodClient) DeletePod(namespace string, name string) (*model.Status, error) {
	status := &model.Status{}
	path := fmt.Sprintf(PodByNamePath, namespace, name)
	err := c.client.doDelete(path, status)
	return status, err
}

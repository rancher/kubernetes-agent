package kubernetesclient

import (
	"fmt"

	"github.com/rancher/kubernetes-model/model"
)

const ReplicationControllerPath string = "/api/v1/namespaces/%s/replicationcontrollers"
const ReplicationControllerByNamePath string = "/api/v1/namespaces/%s/replicationcontrollers/%s"

type ReplicationControllerOperations interface {
	ByName(namespace string, name string) (*model.ReplicationController, error)
	CreateReplicationController(namespace string, resource *model.ReplicationController) (*model.ReplicationController, error)
	ReplaceReplicationController(namespace string, resource *model.ReplicationController) (*model.ReplicationController, error)
	DeleteReplicationController(namespace string, name string) (*model.Status, error)
}

func newReplicationControllerClient(client *Client) *ReplicationControllerClient {
	return &ReplicationControllerClient{
		client: client,
	}
}

type ReplicationControllerClient struct {
	client *Client
}

func (c *ReplicationControllerClient) ByName(namespace string, name string) (*model.ReplicationController, error) {
	resp := &model.ReplicationController{}
	path := fmt.Sprintf(ReplicationControllerByNamePath, namespace, name)
	err := c.client.doGet(path, resp)
	return resp, err
}

func (c *ReplicationControllerClient) CreateReplicationController(namespace string, resource *model.ReplicationController) (*model.ReplicationController, error) {
	resp := &model.ReplicationController{}
	path := fmt.Sprintf(ReplicationControllerPath, namespace)
	err := c.client.doPost(path, resource, resp)
	return resp, err
}

func (c *ReplicationControllerClient) ReplaceReplicationController(namespace string, resource *model.ReplicationController) (*model.ReplicationController, error) {
	resp := &model.ReplicationController{}
	path := fmt.Sprintf(ReplicationControllerByNamePath, namespace, resource.Metadata.Name)
	err := c.client.doPut(path, resource, resp)
	return resp, err
}

func (c *ReplicationControllerClient) DeleteReplicationController(namespace string, name string) (*model.Status, error) {
	status := &model.Status{}
	path := fmt.Sprintf(ReplicationControllerByNamePath, namespace, name)
	err := c.client.doDelete(path, status)
	return status, err
}

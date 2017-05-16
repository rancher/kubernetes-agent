package kubernetesevents

import (
	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
)

func NewChangeHandler(rancherClient *client.RancherClient, kubernetesClient *kubernetesclient.Client, kindHandled string) *ChangeHandler {
	return &ChangeHandler{
		rancherClient: rancherClient,
		kClient:       kubernetesClient,
		kindHandled:   kindHandled,
	}
}

type ChangeHandler struct {
	rancherClient *client.RancherClient
	kClient       *kubernetesclient.Client
	kindHandled   string
}

func (h *ChangeHandler) GetKindHandled() string {
	return h.kindHandled
}

func (h *ChangeHandler) Handle(event model.WatchEvent) error {
	_, err := h.rancherClient.Publish.Create(&client.Publish{
		Name: "service.kubernetes.change",
		Data: map[string]interface{}{
			"type":   event.Type,
			"object": event.Object,
		},
	})

	return err
}

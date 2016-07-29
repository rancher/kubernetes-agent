package kubernetesevents

import (
	"fmt"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

func NewHandler(rancherClient *client.RancherClient, kubernetesClient *kubernetesclient.Client, kindHandled string) (Handler, error) {
	var handler Handler
	if kindHandled == ServiceKind {
		handler = &GenericHandler{
			rancherClient: rancherClient,
			kClient:       kubernetesClient,
			kindHandled:   kindHandled,
		}
	} else if kindHandled == NamespaceKind {
		handler = &NamespaceHandler{
			rancherClient: rancherClient,
			kClient:       kubernetesClient,
			kindHandled:   kindHandled,
		}
	} else {
		return nil, fmt.Errorf("Invalid Handler type %s", kindHandled)
	}
	return handler, nil
}

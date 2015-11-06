package kubernetesclient

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/kubernetes-model/model"
)

func TestService(t *testing.T) {
	client := NewClient("http://localhost:8080", true)

	svcName := "test1"
	err := cleanupService(client, "default", svcName)
	if err != nil {
		t.Fatal(err)
	}

	meta := &model.ObjectMeta{
		Name: svcName,
	}

	selector := map[string]interface{}{
		"foo": "bar",
	}
	ports := make([]model.ServicePort, 0)
	port := model.ServicePort{
		Protocol:   "TCP",
		Port:       8888,
		TargetPort: 8888,
	}
	ports = append(ports, port)
	spec := &model.ServiceSpec{
		Selector:        selector,
		SessionAffinity: "None",
		Ports:           ports,
		Type:            "ClusterIP",
	}
	svc := &model.Service{
		Metadata: meta,
		Spec:     spec,
	}

	respSvc, err := client.Service.CreateService("default", svc)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Created service: %s", respSvc)

	gotService, err := client.Service.ByName("default", svcName)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("Service response: %+v", gotService)

	status, err := client.Service.DeleteService("default", svcName)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("Delete response: %+v", status)
}

func cleanupService(client *Client, namespace string, name string) error {
	_, err := client.Service.DeleteService(namespace, name)
	// _, err := client.Service.ByName(namespace, name)
	if err != nil {
		if apiError, ok := err.(*ApiError); ok && apiError.StatusCode == 404 {
			return nil
		} else {
			return err
		}
	}
	return nil
}

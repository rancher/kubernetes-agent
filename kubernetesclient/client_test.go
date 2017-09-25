package kubernetesclient

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/v1"
)

func TestService(t *testing.T) {
	client := NewClient("http://localhost:8080", true)

	svcName := "test1"
	err := cleanupService(client, "default", svcName)
	if err != nil {
		t.Fatal(err)
	}

	meta := metav1.ObjectMeta{
		Name: svcName,
	}

	selector := map[string]string{
		"foo": "bar",
	}
	ports := make([]v1.ServicePort, 0)
	port := v1.ServicePort{
		Protocol:   "TCP",
		Port:       8888,
		TargetPort: intstr.IntOrString{Type: 0, IntVal: int32(8888)},
	}
	ports = append(ports, port)
	spec := v1.ServiceSpec{
		Selector:        selector,
		SessionAffinity: "None",
		Ports:           ports,
		Type:            "ClusterIP",
	}
	svc := &v1.Service{
		ObjectMeta: meta,
		Spec:       spec,
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

	err = client.Service.DeleteService("default", svcName)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("Service Deleted")
}

func cleanupService(client *Client, namespace string, name string) error {
	err := client.Service.DeleteService(namespace, name)
	if err != nil {
		if k8sErr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

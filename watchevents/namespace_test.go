package watchevents

import (
	"time"

	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"gopkg.in/check.v1"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type NamespacehandlerTestSuite struct {
	kClient *kubernetesclient.Client
	events  chan client.ExternalServiceEvent
}

var _ = check.Suite(&NamespacehandlerTestSuite{})

func (s *NamespacehandlerTestSuite) SetUpSuite(c *check.C) {
	s.events = make(chan client.ExternalServiceEvent, 10)
	s.kClient = kubernetesclient.NewClient(conf.KubernetesURL, true)
	mock := &MockServiceEventOperations{
		events: s.events,
	}
	mockRancherClient := &client.RancherClient{
		ExternalServiceEvent: mock,
	}

	nsHandler := NewNamespaceHandler(mockRancherClient, s.kClient)
	nsHandler.Start()
	defer nsHandler.Stop()
	time.Sleep(time.Second)
}

func (s *NamespacehandlerTestSuite) TestHandler(c *check.C) {
	nsname := "test-ns-1"
	cleanup_ns(s.kClient, "namespace", "test-ns-1", c)

	meta := metav1.ObjectMeta{Name: nsname}
	ns := &v1.Namespace{
		ObjectMeta: meta,
	}

	respNs, err := s.kClient.Namespace.CreateNamespace(ns)
	if err != nil {
		c.Fatal(err)
	}

	err = s.kClient.Namespace.DeleteNamespace(nsname)

	var gotDelete bool
	for !gotDelete {
		select {
		case event := <-s.events:
			c.Logf("%#v %+v", event, event)
			svc := event.Service
			service := svc.(client.Service)
			c.Logf("EXPECTED %s; EVENT %s", string(respNs.UID), event)
			if event.EventType == "stack.remove" {
				c.Assert(service.Kind, check.Equals, "kubernetesService")
				c.Assert(event.ExternalId, check.Equals, "kubernetes://"+string(respNs.UID))
				gotDelete = true
			}
		case <-time.After(time.Second * 5):
			c.Fatalf("Timed out waiting for event.")

		}
	}
}

func cleanup_ns(client *kubernetesclient.Client, resourceType string, namespace string, c *check.C) error {
	var err error
	switch resourceType {
	case "namespace":
		err = client.Namespace.DeleteNamespace(namespace)
	default:
		c.Fatalf("Unknown type for cleanup: %s", resourceType)
	}
	if err != nil {
		if k8sErr.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

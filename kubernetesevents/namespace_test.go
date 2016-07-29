package kubernetesevents

import (
	"gopkg.in/check.v1"
	"time"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-model/model"

	"github.com/rancher/kubernetes-agent/kubernetesclient"
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

	nsHandler, _ := NewHandler(mockRancherClient, s.kClient, NamespaceKind)
	handlers := []Handler{nsHandler}
	go ConnectToEventStream(handlers, conf)
	time.Sleep(time.Second)
}

func (s *NamespacehandlerTestSuite) TestHandler(c *check.C) {
	nsname := "test-ns-1"
	cleanup_ns(s.kClient, "namespace", "test-ns-1", c)

	meta := &model.ObjectMeta{Name: nsname}
	ns := &model.Namespace{
		Metadata: meta,
	}

	respNs, err := s.kClient.Namespace.CreateNamespace(ns)
	if err != nil {
		c.Fatal(err)
	}

	_, err = s.kClient.Namespace.DeleteNamespace(nsname)

	var gotDelete bool
	for !gotDelete {
		select {
		case event := <-s.events:
			svc := event.Service
			service := svc.(client.Service)
			c.Logf("EXPECTED %s; EVENT %s", respNs.Metadata.Uid, event)
			if event.EventType == "stack.remove" {
				c.Assert(service.Kind, check.Equals, "kubernetesService")
				c.Assert(event.ExternalId, check.Equals, "kubernetes://"+respNs.Metadata.Uid)
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
		_, err = client.Namespace.DeleteNamespace(namespace)
	default:
		c.Fatalf("Unknown type for cleanup: %s", resourceType)
	}
	if err != nil {
		if apiError, ok := err.(*kubernetesclient.ApiError); ok && apiError.StatusCode == 404 {
			return nil
		} else {
			return err
		}
	}
	return nil
}

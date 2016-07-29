package kubernetesevents

import (
	"gopkg.in/check.v1"
	"testing"
	"time"

	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-model/model"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

var conf = config.Config{
	KubernetesURL:   "http://localhost:8080",
	CattleURL:       "http://localhost:8082",
	CattleAccessKey: "agent",
	CattleSecretKey: "agentpass",
	WorkerCount:     10,
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	check.TestingT(t)
}

type GenerichandlerTestSuite struct {
	kClient *kubernetesclient.Client
	events  chan client.ExternalServiceEvent
}

var _ = check.Suite(&GenerichandlerTestSuite{})

func (s *GenerichandlerTestSuite) SetUpSuite(c *check.C) {
	s.events = make(chan client.ExternalServiceEvent, 10)
	s.kClient = kubernetesclient.NewClient(conf.KubernetesURL, true)
	mock := &MockServiceEventOperations{
		events: s.events,
	}
	mockRancherClient := &client.RancherClient{
		ExternalServiceEvent: mock,
	}

	svcHandler, _ := NewHandler(mockRancherClient, s.kClient, ServiceKind)
	handlers := []Handler{svcHandler}
	go ConnectToEventStream(handlers, conf)
	time.Sleep(time.Second)
}

func (s *GenerichandlerTestSuite) TestService(c *check.C) {
	svcName := "test-service-1"
	cleanup(s.kClient, "service", "default", svcName, c)

	meta := &model.ObjectMeta{Name: svcName}
	selector := map[string]interface{}{"foo": "bar", "env": "dev"}
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

	respSvc, err := s.kClient.Service.CreateService("default", svc)
	if err != nil {
		c.Fatal(err)
	}

	newSelector := map[string]interface{}{"env": "prod"}
	respSvc.Spec.Selector = newSelector
	_, err = s.kClient.Service.ReplaceService("default", respSvc)
	if err != nil {
		c.Fatal(err)
	}

	_, err = s.kClient.Service.DeleteService("default", svcName)
	if err != nil {
		c.Fatal(err)
	}

	var gotCreate, gotMod, gotDelete bool
	for !gotCreate || !gotMod || !gotDelete {
		select {
		case event := <-s.events:
			svc := event.Service
			service := svc.(client.Service)
			c.Logf("EXPECTED %s; EVENT %s", respSvc.Metadata.Uid, event)
			if event.ExternalId == respSvc.Metadata.Uid {
				if event.EventType == "service.create" {
					c.Assert(service.Kind, check.Equals, "kubernetesService")
					c.Assert(service.Name, check.Equals, svcName)
					c.Assert(service.ExternalId, check.Equals, respSvc.Metadata.Uid)
					c.Assert(service.SelectorContainer, check.Matches, "foo=bar,env=dev,io.kubernetes.pod.namespace=default|env=dev,foo=bar,io.kubernetes.pod.namespace=default|io.kubernetes.pod.namespace=default,foo=bar,env=dev|io.kubernetes.pod.namespace=default,env=dev,foo=bar|foo=bar,io.kubernetes.pod.namespace=default,env=dev|env=dev,io.kubernetes.pod.namespace=default,foo=bar")

					env := event.Environment.(map[string]string)
					c.Assert(env["name"], check.Equals, "default")
					kEnv, err := s.kClient.Namespace.ByName("default")
					if err != nil {
						c.Fatal(err)
					}
					c.Assert(env["externalId"], check.Equals, "kubernetes://"+kEnv.Metadata.Uid)
					gotCreate = true
				} else if event.EventType == "service.update" {
					c.Assert(service.Kind, check.Equals, "kubernetesService")
					c.Assert(service.Name, check.Equals, svcName)
					c.Assert(service.ExternalId, check.Equals, respSvc.Metadata.Uid)
					c.Assert(service.SelectorContainer, check.Matches, "env=prod,io.kubernetes.pod.namespace=default|io.kubernetes.pod.namespace=default,env=prod")
					gotMod = true
				} else if event.EventType == "service.remove" {
					gotDelete = true
				}
			}
		case <-time.After(time.Second * 5):
			c.Fatalf("Timed out waiting for event.")

		}
	}
}

type MockServiceEventOperations struct {
	client.ExternalServiceEventClient
	events chan<- client.ExternalServiceEvent
}

func (m *MockServiceEventOperations) Create(event *client.ExternalServiceEvent) (*client.ExternalServiceEvent, error) {
	m.events <- *event
	return nil, nil
}

func cleanup(client *kubernetesclient.Client, resourceType string, namespace string, name string, c *check.C) error {
	var err error
	switch resourceType {
	case "service":
		_, err = client.Service.DeleteService(namespace, name)
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

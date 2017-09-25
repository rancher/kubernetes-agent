package rancherevents

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/check.v1"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"

	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v2"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/dockerclient"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-agent/rancherevents/eventhandlers"
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

type ListenerTestSuite struct {
	kClient     *kubernetesclient.Client
	dClient     *docker.Client
	publishChan chan client.Publish
	mockRClient *client.RancherClient
}

var _ = check.Suite(&ListenerTestSuite{})

func (s *ListenerTestSuite) SetUpSuite(c *check.C) {
	s.publishChan = make(chan client.Publish, 10)

	s.kClient = kubernetesclient.NewClient(conf.KubernetesURL, true)

	mock := &MockPublishOperations{
		publishChan: s.publishChan,
	}
	s.mockRClient = &client.RancherClient{
		Publish: mock,
	}

	dClient, err := dockerclient.NewDockerClient()
	if err != nil {
		c.Fatal(err)
	}
	s.dClient = dClient
}

func (s *ListenerTestSuite) TestSyncHandler(c *check.C) {
	pod, containers, err := s.createPod(c)
	if err != nil {
		c.Fatal(err)
	}

	log.Info(containers)

	container := containers[0]
	event := &revents.Event{
		ReplyTo: "event-1",
		ID:      "event-id-1",
		Data: map[string]interface{}{
			"instanceHostMap": map[string]interface{}{
				"instance": map[string]interface{}{
					"externalId": container.ID,
					"data": map[string]interface{}{
						"fields": map[string]interface{}{
							"labels": map[string]interface{}{
								"io.kubernetes.pod.namespace":  "default",
								"io.kubernetes.pod.name":       "pod-test-1",
								"io.kubernetes.container.name": "POD",
								"io.kubernetes.pod.uid":        string(pod.ObjectMeta.UID),
							},
						},
					},
				},
			},
		},
	}
	sh := eventhandlers.NewProvideLablesHandler(s.kClient)

	err = sh.Handler(event, s.mockRClient)
	if err != nil {
		c.Fatal(err)
	}

	pub := <-s.publishChan
	c.Assert(pub.Name, check.Equals, "event-1")
	c.Assert(pub.PreviousIds, check.DeepEquals, []string{"event-id-1"})
	instance := get(pub.Data, "instance")
	data := get(instance, "+data")
	fields := get(data, "+fields")
	newLabels := get(fields, "+labels")
	c.Assert(newLabels, check.DeepEquals,
		map[string]string{
			"env": "dev",
			"io.rancher.service.deployment.unit": string(pod.ObjectMeta.UID),
			"io.rancher.stack.name":              "default",
			"io.rancher.container.display_name":  "pod-test-1",
			"io.rancher.container.network":       "true",
			"io.rancher.service.launch.config":   "io.rancher.service.primary.launch.config",
		})
}

func get(theMap interface{}, key string) interface{} {
	if castedMap, ok := theMap.(map[string]interface{}); ok {
		return castedMap[key]
	} else {
		return nil
	}
}

func (s *ListenerTestSuite) createPod(c *check.C) (*v1.Pod, []docker.APIContainers, error) {
	podName := "pod-test-1"
	cleanup(s.kClient, "pod", "default", podName, c)

	podLabels := map[string]string{"env": "dev"}

	podMeta := metav1.ObjectMeta{
		Labels: podLabels,
		Name:   podName,
	}
	container := v1.Container{
		Name:            "pod-test",
		Image:           "nginx",
		ImagePullPolicy: "IfNotPresent",
	}
	containers := []v1.Container{container}
	podSpec := v1.PodSpec{
		Containers:    containers,
		RestartPolicy: "Always",
		DNSPolicy:     "ClusterFirst",
	}

	pod := &v1.Pod{
		ObjectMeta: podMeta,
		Spec:       podSpec,
	}

	pod, err := s.kClient.Pod.CreatePod("default", pod)
	if err != nil {
		c.Fatal(err)
	}

	opts := docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {"io.kubernetes.pod.name=pod-test-1"},
		},
	}
	for i := 0; i < 10; i++ {
		containers, err := s.dClient.ListContainers(opts)
		if err != nil {
			return nil, nil, err
		}
		if len(containers) > 0 {
			return pod, containers, nil
		}
		<-time.After(time.Second * 5)
	}
	c.Fatal("Timed out waiting for containers to be created for pod.")
	return nil, nil, nil
}

type MockPublishOperations struct {
	client.PublishClient
	publishChan chan<- client.Publish
}

func (m *MockPublishOperations) Create(publish *client.Publish) (*client.Publish, error) {
	m.publishChan <- *publish
	return nil, nil
}

func cleanup(client *kubernetesclient.Client, resourceType string, namespace string, name string, c *check.C) error {
	var err error
	switch resourceType {
	case "pod":
		err = client.Pod.DeletePod(namespace, name)
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

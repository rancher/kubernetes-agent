package rancherevents

import (
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"gopkg.in/check.v1"

	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-model/model"

	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/dockerclient"
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
	containers, err := s.createPod(c)
	if err != nil {
		c.Fatal(err)
	}

	log.Info(containers)

	container := containers[0]
	event := &revents.Event{
		ReplyTo: "event-1",
		Id:      "event-id-1",
		Data: map[string]interface{}{
			"instanceHostMap": map[string]interface{}{
				"instance": map[string]interface{}{
					"externalId": container.ID,
				},
			},
		},
	}
	sh := syncHandler{
		kClient: s.kClient,
		dClient: s.dClient,
	}

	err = sh.Handler(event, s.mockRClient)
	if err != nil {
		c.Fatal(err)
	}

	pub := <-s.publishChan
	c.Assert(pub.Name, check.Equals, "event-1")
	c.Assert(pub.PreviousIds, check.DeepEquals, []string{"event-id-1"})
	instance := get(pub.Data["instanceHostMap"], "instance")
	data := get(instance, "+data")
	fields := get(data, "+fields")
	newLabels := get(fields, "+labels")
	c.Assert(newLabels, check.DeepEquals, map[string]string{"env": "dev"})
}

func get(theMap interface{}, key string) interface{} {
	if castedMap, ok := theMap.(map[string]interface{}); ok {
		return castedMap[key]
	} else {
		return nil
	}
}

func (s *ListenerTestSuite) createPod(c *check.C) ([]docker.APIContainers, error) {
	podName := "pod-test-1"
	cleanup(s.kClient, "pod", "default", podName, c)

	podLabels := map[string]interface{}{"env": "dev"}
	podMeta := &model.ObjectMeta{
		Labels: podLabels,
		Name:   podName,
	}
	container := model.Container{
		Name:            "pod-test",
		Image:           "nginx",
		ImagePullPolicy: "IfNotPresent",
	}
	containers := []model.Container{container}
	podSpec := &model.PodSpec{
		Containers:    containers,
		RestartPolicy: "Always",
		DnsPolicy:     "ClusterFirst",
	}

	pod := &model.Pod{
		Metadata: podMeta,
		Spec:     podSpec,
	}

	_, err := s.kClient.Pod.CreatePod("default", pod)
	if err != nil {
		c.Fatal(err)
	}

	opts := docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {"io.kubernetes.pod.name=default/pod-test-1"},
		},
	}
	for i := 0; i < 10; i++ {
		containers, err := s.dClient.ListContainers(opts)
		if err != nil {
			return nil, err
		}
		if len(containers) > 0 {
			return containers, nil
		}
		<-time.After(time.Second * 5)
	}
	c.Fatal("Timed out waiting for containers to be created for pod.")
	return nil, nil
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
		_, err = client.Pod.DeletePod(namespace, name)
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

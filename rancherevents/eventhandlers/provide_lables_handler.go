package eventhandlers

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v2"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	util "github.com/rancher/kubernetes-agent/rancherevents/util"
	k8sErr "k8s.io/apimachinery/pkg/api/errors"
)

type syncHandler struct {
	kClient *kubernetesclient.Client
}

func NewProvideLablesHandler(kClient *kubernetesclient.Client) *syncHandler {
	return &syncHandler{
		kClient: kClient,
	}
}

func (h *syncHandler) Handler(event *events.Event, cli *client.RancherClient) error {
	log := logrus.WithFields(logrus.Fields{
		"eventName":  event.Name,
		"eventID":    event.ID,
		"resourceID": event.ResourceID,
	})
	log.Infof("Rancher event: %#v", event)
	labels := map[string]string{}

	containerLabels, err := h.parseContainerLabels(event)
	if err != nil {
		log.Errorf("Failed to read labels", err)
		return util.CreateAndPublishReply(event, cli)
	}

	namespace := containerLabels["io.kubernetes.pod.namespace"]
	name := containerLabels["io.kubernetes.pod.name"]
	if namespace == "" || name == "" {
		return util.CreateAndPublishReply(event, cli)
	}

	labels["io.rancher.service.deployment.unit"] = containerLabels["io.kubernetes.pod.uid"]
	labels["io.rancher.stack.name"] = namespace

	if isPodContainer(containerLabels) {
		if !isHostNetwork(event) {
			labels["io.rancher.container.network"] = "true"
		}
		labels["io.rancher.service.launch.config"] = "io.rancher.service.primary.launch.config"
		labels["io.rancher.container.display_name"] = containerLabels["io.kubernetes.pod.name"]
		if found, err := h.copyPodLabels(namespace, name, labels); err != nil {
			return util.ErrorReply(event, cli, err)
		} else if !found {
			return util.CreateAndPublishReply(event, cli)
		}
	} else {
		labels["io.rancher.container.display_name"] = containerLabels["io.kubernetes.container.name"]
	}

	return h.replyWithLabels(event, cli, labels)
}

func isPodContainer(containerLabels map[string]string) bool {
	return containerLabels["io.kubernetes.container.name"] == "POD"
}

func (h *syncHandler) copyPodLabels(namespace, name string, labels map[string]string) (bool, error) {
	pod, err := h.kClient.Pod.ByName(namespace, name)
	if err != nil {
		if k8sErr.IsNotFound(err) {
			return false, nil
		}
		return true, errors.Wrap(err, "lookup pod")
	}

	for key, v := range pod.GetLabels() {
		labels[key] = v
	}

	return true, nil
}

func (h *syncHandler) replyWithLabels(event *events.Event, cli *client.RancherClient, labels map[string]string) error {
	reply := util.NewReply(event)
	reply.ResourceType = event.ResourceType
	reply.ResourceId = event.ResourceID
	reply.Data = map[string]interface{}{
		"instance": map[string]interface{}{
			"+data": map[string]interface{}{
				"+fields": map[string]interface{}{
					"+labels": labels,
				},
			},
		},
	}
	logrus.WithField("eventID", event.ID).Infof("Reply: %+v", reply)
	return util.PublishReply(reply, cli)
}

func (h *syncHandler) parseContainerLabels(event *events.Event) (map[string]string, error) {
	labels := GetStringMap(event.Data, "instanceHostMap", "instance", "data", "fields", "labels")
	if len(labels) == 0 {
		labels = GetStringMap(event.Data, "instance", "data", "fields", "labels")
	}

	return labels, nil
}

func isHostNetwork(event *events.Event) bool {
	data := event.Data["instanceHostMap"]
	var instance util.InstanceHostMapData
	mapstructure.Decode(data, &instance)
	networkMode := instance.Instance.Data.DockerContainer.HostConfig.NetworkMode
	if networkMode == "" {
		logrus.Warnf("Couldn't decode %+v to instanceData.", data)
		return false
	}

	return strings.EqualFold("host", networkMode)
}

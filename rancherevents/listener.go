package rancherevents

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"strings"

	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/config"
	"github.com/rancher/kubernetes-agent/dockerclient"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

func ConnectToEventStream(conf config.Config) error {

	kClient := kubernetesclient.NewClient(conf.KubernetesURL, false)
	dClient, err := dockerclient.NewDockerClient()
	if err != nil {
		return err
	}
	sh := syncHandler{
		kClient: kClient,
		dClient: dClient,
	}
	ph := PingHandler{}

	eventHandlers := map[string]revents.EventHandler{
		"compute.instance.providelabels": sh.Handler,
		"config.update":                  ph.Handler,
		"ping":                           ph.Handler,
	}

	router, err := revents.NewEventRouter("", 0, conf.CattleURL, conf.CattleAccessKey, conf.CattleSecretKey, nil, eventHandlers, "", conf.WorkerCount)
	if err != nil {
		return err
	}
	err = router.StartWithoutCreate(nil)
	return err
}

type syncHandler struct {
	kClient *kubernetesclient.Client
	dClient *docker.Client
}

func (h *syncHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	log.Infof("Received event: %s, %s", event.Name, event.Id)

	namespace, name := h.getPod(event)
	if namespace == "" || name == "" {
		reply := newReply(event)
		reply.ResourceType = event.ResourceType
		reply.ResourceId = event.ResourceId
		err := publishReply(reply, cli)
		if err != nil {
			return err
		}
		return nil
	}

	pod, err := h.kClient.Pod.ByName(namespace, name)
	if err != nil {
		if apiErr, ok := err.(*client.ApiError); ok && apiErr.StatusCode == 404 {
			reply := newReply(event)
			reply.ResourceType = event.ResourceType
			reply.ResourceId = event.ResourceId
			err := publishReply(reply, cli)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	/*
		  Building this:
		  {instancehostmap: {
			  instance: {
				  +data: {
					  +fields: {
						labels: {
	*/
	replyData := make(map[string]interface{})
	ihm := make(map[string]interface{})
	i := make(map[string]interface{})
	data := make(map[string]interface{})
	fields := make(map[string]interface{})
	labels := make(map[string]string)

	replyData["instanceHostMap"] = ihm
	ihm["instance"] = i
	i["+data"] = data
	data["+fields"] = fields
	fields["+labels"] = labels

	for key, v := range pod.Metadata.Labels {
		if val, ok := v.(string); ok {
			labels[key] = val
		}
	}

	labels["io.rancher.service.deployment.unit"] = pod.Metadata.Uid
	labels["io.rancher.stack.name"] = pod.Metadata.Namespace

	reply := newReply(event)
	reply.ResourceType = "instanceHostMap"
	reply.ResourceId = event.ResourceId
	reply.Data = replyData
	log.Infof("Reply: %+v", reply)
	err = publishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

func (h *syncHandler) getPod(event *revents.Event) (string, string) {
	// TODO Rewrite this horror
	data := event.Data
	if ihm, ok := data["instanceHostMap"]; ok {
		if ihmMap, ok := ihm.(map[string]interface{}); ok {
			if i, ok := ihmMap["instance"]; ok {
				if iMap, ok := i.(map[string]interface{}); ok {
					if eId, ok := iMap["externalId"]; ok {
						if externalId, ok := eId.(string); ok {
							container, err := h.dClient.InspectContainer(externalId)
							if err == nil {
								if podName, ok := container.Config.Labels["io.kubernetes.pod.name"]; ok {
									parts := strings.SplitN(podName, "/", 2)
									if len(parts) == 2 {
										return parts[0], parts[1]
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return "", ""
}

type PingHandler struct {
}

func (h *PingHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	reply := newReply(event)
	if reply.Name == "" {
		return nil
	}
	err := publishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

func newReply(event *revents.Event) *client.Publish {
	return &client.Publish{
		Name:        event.ReplyTo,
		PreviousIds: []string{event.Id},
	}
}

func publishReply(reply *client.Publish, apiClient *client.RancherClient) error {
	_, err := apiClient.Publish.Create(reply)
	return err
}

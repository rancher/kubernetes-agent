package eventhandlers

import (
	log "github.com/Sirupsen/logrus"
	"strings"

	revents "github.com/rancher/go-machine-service/events"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	util "github.com/rancher/kubernetes-agent/rancherevents/util"
)

type syncHandler struct {
	kClient *kubernetesclient.Client
}

func NewProvideLablesHandler(kClient *kubernetesclient.Client) *syncHandler {
	return &syncHandler{
		kClient: kClient,
	}
}

func (h *syncHandler) Handler(event *revents.Event, cli *client.RancherClient) error {
	log.Infof("Received event: Name: %s, Event Id: %s, Resource Id: %s", event.Name, event.Id, event.ResourceId)

	namespace, name := h.getPod(event)
	if namespace == "" || name == "" {
		if err := util.CreateAndPublishReply(event, cli); err != nil {
			return err
		}
		return nil
	}

	pod, err := h.kClient.Pod.ByName(namespace, name)
	if err != nil {
		log.Warnf("Error looking up pod: %#v", err)
		if apiErr, ok := err.(*client.ApiError); ok && apiErr.StatusCode == 404 {
			if err := util.CreateAndPublishReply(event, cli); err != nil {
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

	reply := util.NewReply(event)
	reply.ResourceType = "instanceHostMap"
	reply.ResourceId = event.ResourceId
	reply.Data = replyData
	log.Infof("Reply: %+v", reply)
	err = util.PublishReply(reply, cli)
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
					if d, ok := iMap["data"]; ok {
						if dMap, ok := d.(map[string]interface{}); ok {
							if f, ok := dMap["fields"]; ok {
								if fMap, ok := f.(map[string]interface{}); ok {
									if labels, ok := fMap["labels"]; ok {
										if lMap, ok := labels.(map[string]interface{}); ok {
											if l, ok := lMap["io.kubernetes.pod.name"]; ok {
												if label, ok := l.(string); ok {
													parts := strings.SplitN(label, "/", 2)
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
				}
			}
		}
	}

	return "", ""
}

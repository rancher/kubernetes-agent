package eventhandlers

import (
	revents "github.com/rancher/event-subscriber/events"
	"github.com/rancher/go-rancher/v2"
)

func NewReply(event *revents.Event) *client.Publish {
	return &client.Publish{
		Name:         event.ReplyTo,
		PreviousIds:  []string{event.ID},
		ResourceType: event.ResourceType,
		ResourceId:   event.ResourceID,
	}
}

func PublishReply(reply *client.Publish, apiClient *client.RancherClient) error {
	_, err := apiClient.Publish.Create(reply)
	return err
}

func CreateAndPublishReply(event *revents.Event, cli *client.RancherClient) error {
	reply := NewReply(event)
	if reply.Name == "" {
		return nil
	}
	err := PublishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

func ErrorReply(event *revents.Event, cli *client.RancherClient, eventError error) error {
	reply := NewReply(event)
	if reply.Name == "" {
		return nil
	}
	reply.Transitioning = "error"
	reply.TransitioningMessage = eventError.Error()
	err := PublishReply(reply, cli)
	if err != nil {
		return err
	}
	return nil
}

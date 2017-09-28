package events

type Event struct {
	Name                 string                 `json:"name,omitempty"`
	ID                   string                 `json:"id,omitempty"`
	PreviousIds          string                 `json:"previousIds,omitempty"`
	ReplyTo              string                 `json:"replyTo,omitempty"`
	ResourceID           string                 `json:"resourceId,omitempty"`
	ResourceType         string                 `json:"resourceType,omitempty"`
	Transitioning        string                 `json:"transitioning,omitempty"`
	TransitioningMessage string                 `json:"transitioningMessage,omitempty"`
	Data                 map[string]interface{} `json:"data,omitempty"`
	Time                 float64                `json:"time,omitempty"`
}

type ReplyEvent struct {
	Name        string                 `json:"name"`
	PreviousIds []string               `json:"previousIds"`
	Data        map[string]interface{} `json:"data"`
}

func NewReplyEvent(replyTo string, eventID string) *ReplyEvent {
	return &ReplyEvent{Name: replyTo, PreviousIds: []string{eventID}}
}

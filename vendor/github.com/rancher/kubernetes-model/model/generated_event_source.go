package model

const (
	EVENT_SOURCE_TYPE = "v1.EventSource"
)

type EventSource struct {
	Component string `json:"component,omitempty" yaml:"component,omitempty"`

	Host string `json:"host,omitempty" yaml:"host,omitempty"`
}

package model

const (
	WATCH_EVENT_TYPE = "json.WatchEvent"
)

type WatchEvent struct {
	Object interface{} `json:"object,omitempty" yaml:"object,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

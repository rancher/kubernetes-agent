package model

const (
	EVENT_TYPE = "v1.Event"
)

type Event struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Count int32 `json:"count,omitempty" yaml:"count,omitempty"`

	FirstTimestamp string `json:"firstTimestamp,omitempty" yaml:"first_timestamp,omitempty"`

	InvolvedObject *ObjectReference `json:"involvedObject,omitempty" yaml:"involved_object,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	LastTimestamp string `json:"lastTimestamp,omitempty" yaml:"last_timestamp,omitempty"`

	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`

	Source *EventSource `json:"source,omitempty" yaml:"source,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

package model

const (
	STATUS_CAUSE_TYPE = "unversioned.StatusCause"
)

type StatusCause struct {
	Field string `json:"field,omitempty" yaml:"field,omitempty"`

	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

package model

const (
	CAPABILITIES_TYPE = "v1.Capabilities"
)

type Capabilities struct {
	Add []Capability `json:"add,omitempty" yaml:"add,omitempty"`

	Drop []Capability `json:"drop,omitempty" yaml:"drop,omitempty"`
}

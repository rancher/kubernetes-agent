package model

const (
	POD_CONDITION_TYPE = "v1.PodCondition"
)

type PodCondition struct {
	LastProbeTime string `json:"lastProbeTime,omitempty" yaml:"last_probe_time,omitempty"`

	LastTransitionTime string `json:"lastTransitionTime,omitempty" yaml:"last_transition_time,omitempty"`

	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`

	Status string `json:"status,omitempty" yaml:"status,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

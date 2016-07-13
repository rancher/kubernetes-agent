package model

const (
	POD_TYPE = "v1.Pod"
)

type Pod struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *PodSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *PodStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

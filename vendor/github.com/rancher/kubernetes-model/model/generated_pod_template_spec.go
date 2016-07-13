package model

const (
	POD_TEMPLATE_SPEC_TYPE = "v1.PodTemplateSpec"
)

type PodTemplateSpec struct {
	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *PodSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

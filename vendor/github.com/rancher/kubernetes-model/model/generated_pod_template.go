package model

const (
	POD_TEMPLATE_TYPE = "v1.PodTemplate"
)

type PodTemplate struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Template *PodTemplateSpec `json:"template,omitempty" yaml:"template,omitempty"`
}

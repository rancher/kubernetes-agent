package model

const (
	SERVICE_TYPE = "v1.Service"
)

type Service struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *ServiceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *ServiceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

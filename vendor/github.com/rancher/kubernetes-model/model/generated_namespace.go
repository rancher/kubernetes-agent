package model

const (
	NAMESPACE_TYPE = "v1.Namespace"
)

type Namespace struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *NamespaceSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *NamespaceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

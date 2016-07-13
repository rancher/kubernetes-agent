package model

const (
	BINDING_TYPE = "v1.Binding"
)

type Binding struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Target *ObjectReference `json:"target,omitempty" yaml:"target,omitempty"`
}

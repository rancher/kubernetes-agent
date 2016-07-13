package model

const (
	NODE_TYPE = "v1.Node"
)

type Node struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *NodeSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *NodeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

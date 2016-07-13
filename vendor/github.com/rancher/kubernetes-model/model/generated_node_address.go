package model

const (
	NODE_ADDRESS_TYPE = "v1.NodeAddress"
)

type NodeAddress struct {
	Address string `json:"address,omitempty" yaml:"address,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

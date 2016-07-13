package model

const (
	NODE_SPEC_TYPE = "v1.NodeSpec"
)

type NodeSpec struct {
	ExternalID string `json:"externalID,omitempty" yaml:"external_id,omitempty"`

	PodCIDR string `json:"podCIDR,omitempty" yaml:"pod_cidr,omitempty"`

	ProviderID string `json:"providerID,omitempty" yaml:"provider_id,omitempty"`

	Unschedulable bool `json:"unschedulable,omitempty" yaml:"unschedulable,omitempty"`
}

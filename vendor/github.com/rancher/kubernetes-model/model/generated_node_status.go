package model

const (
	NODE_STATUS_TYPE = "v1.NodeStatus"
)

type NodeStatus struct {
	Addresses []NodeAddress `json:"addresses,omitempty" yaml:"addresses,omitempty"`

	Allocatable map[string]interface{} `json:"allocatable,omitempty" yaml:"allocatable,omitempty"`

	Capacity map[string]interface{} `json:"capacity,omitempty" yaml:"capacity,omitempty"`

	Conditions []NodeCondition `json:"conditions,omitempty" yaml:"conditions,omitempty"`

	DaemonEndpoints *NodeDaemonEndpoints `json:"daemonEndpoints,omitempty" yaml:"daemon_endpoints,omitempty"`

	Images []ContainerImage `json:"images,omitempty" yaml:"images,omitempty"`

	NodeInfo *NodeSystemInfo `json:"nodeInfo,omitempty" yaml:"node_info,omitempty"`

	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`
}

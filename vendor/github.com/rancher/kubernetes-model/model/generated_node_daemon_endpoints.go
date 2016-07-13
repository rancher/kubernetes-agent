package model

const (
	NODE_DAEMON_ENDPOINTS_TYPE = "v1.NodeDaemonEndpoints"
)

type NodeDaemonEndpoints struct {
	KubeletEndpoint *DaemonEndpoint `json:"kubeletEndpoint,omitempty" yaml:"kubelet_endpoint,omitempty"`
}

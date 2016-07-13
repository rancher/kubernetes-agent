package model

const (
	NAMESPACE_STATUS_TYPE = "v1.NamespaceStatus"
)

type NamespaceStatus struct {
	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`
}

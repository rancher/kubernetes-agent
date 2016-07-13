package model

const (
	NAMESPACE_SPEC_TYPE = "v1.NamespaceSpec"
)

type NamespaceSpec struct {
	Finalizers []FinalizerName `json:"finalizers,omitempty" yaml:"finalizers,omitempty"`
}

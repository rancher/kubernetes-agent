package model

const (
	LOCAL_OBJECT_REFERENCE_TYPE = "v1.LocalObjectReference"
)

type LocalObjectReference struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

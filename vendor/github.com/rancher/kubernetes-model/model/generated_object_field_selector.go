package model

const (
	OBJECT_FIELD_SELECTOR_TYPE = "v1.ObjectFieldSelector"
)

type ObjectFieldSelector struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	FieldPath string `json:"fieldPath,omitempty" yaml:"field_path,omitempty"`
}

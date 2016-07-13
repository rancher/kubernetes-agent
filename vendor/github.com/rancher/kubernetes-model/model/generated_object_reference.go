package model

const (
	OBJECT_REFERENCE_TYPE = "v1.ObjectReference"
)

type ObjectReference struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	FieldPath string `json:"fieldPath,omitempty" yaml:"field_path,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	ResourceVersion string `json:"resourceVersion,omitempty" yaml:"resource_version,omitempty"`

	Uid string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

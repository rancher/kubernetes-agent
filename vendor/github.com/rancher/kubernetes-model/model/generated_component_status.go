package model

const (
	COMPONENT_STATUS_TYPE = "v1.ComponentStatus"
)

type ComponentStatus struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Conditions []ComponentCondition `json:"conditions,omitempty" yaml:"conditions,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

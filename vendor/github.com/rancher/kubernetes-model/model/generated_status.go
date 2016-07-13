package model

const (
	STATUS_TYPE = "unversioned.Status"
)

type Status struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Code int32 `json:"code,omitempty" yaml:"code,omitempty"`

	Details *StatusDetails `json:"details,omitempty" yaml:"details,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`

	Status string `json:"status,omitempty" yaml:"status,omitempty"`
}

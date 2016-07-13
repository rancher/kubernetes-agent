package model

const (
	LIMIT_RANGE_TYPE = "v1.LimitRange"
)

type LimitRange struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *LimitRangeSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

package model

const (
	LIMIT_RANGE_SPEC_TYPE = "v1.LimitRangeSpec"
)

type LimitRangeSpec struct {
	Limits []LimitRangeItem `json:"limits,omitempty" yaml:"limits,omitempty"`
}

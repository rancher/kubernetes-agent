package model

const (
	LIFECYCLE_TYPE = "v1.Lifecycle"
)

type Lifecycle struct {
	PostStart *Handler `json:"postStart,omitempty" yaml:"post_start,omitempty"`

	PreStop *Handler `json:"preStop,omitempty" yaml:"pre_stop,omitempty"`
}

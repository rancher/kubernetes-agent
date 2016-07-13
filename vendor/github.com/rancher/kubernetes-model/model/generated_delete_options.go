package model

const (
	DELETE_OPTIONS_TYPE = "v1.DeleteOptions"
)

type DeleteOptions struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	GracePeriodSeconds int64 `json:"gracePeriodSeconds,omitempty" yaml:"grace_period_seconds,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

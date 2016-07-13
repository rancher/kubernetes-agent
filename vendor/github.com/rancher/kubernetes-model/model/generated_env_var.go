package model

const (
	ENV_VAR_TYPE = "v1.EnvVar"
)

type EnvVar struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Value string `json:"value,omitempty" yaml:"value,omitempty"`

	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" yaml:"value_from,omitempty"`
}

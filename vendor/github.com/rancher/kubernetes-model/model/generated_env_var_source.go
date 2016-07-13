package model

const (
	ENV_VAR_SOURCE_TYPE = "v1.EnvVarSource"
)

type EnvVarSource struct {
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty" yaml:"config_map_key_ref,omitempty"`

	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" yaml:"field_ref,omitempty"`

	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty" yaml:"secret_key_ref,omitempty"`
}

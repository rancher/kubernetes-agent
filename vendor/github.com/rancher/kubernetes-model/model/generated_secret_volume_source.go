package model

const (
	SECRET_VOLUME_SOURCE_TYPE = "v1.SecretVolumeSource"
)

type SecretVolumeSource struct {
	SecretName string `json:"secretName,omitempty" yaml:"secret_name,omitempty"`
}

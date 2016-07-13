package model

const (
	AZURE_FILE_VOLUME_SOURCE_TYPE = "v1.AzureFileVolumeSource"
)

type AzureFileVolumeSource struct {
	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	SecretName string `json:"secretName,omitempty" yaml:"secret_name,omitempty"`

	ShareName string `json:"shareName,omitempty" yaml:"share_name,omitempty"`
}

package model

const (
	FLEX_VOLUME_SOURCE_TYPE = "v1.FlexVolumeSource"
)

type FlexVolumeSource struct {
	Driver string `json:"driver,omitempty" yaml:"driver,omitempty"`

	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	SecretRef *LocalObjectReference `json:"secretRef,omitempty" yaml:"secret_ref,omitempty"`
}

package model

const (
	RBDVOLUME_SOURCE_TYPE = "v1.RBDVolumeSource"
)

type RBDVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	Keyring string `json:"keyring,omitempty" yaml:"keyring,omitempty"`

	Monitors []string `json:"monitors,omitempty" yaml:"monitors,omitempty"`

	Pool string `json:"pool,omitempty" yaml:"pool,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	SecretRef *LocalObjectReference `json:"secretRef,omitempty" yaml:"secret_ref,omitempty"`

	User string `json:"user,omitempty" yaml:"user,omitempty"`
}

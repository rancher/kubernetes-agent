package model

const (
	VOLUME_MOUNT_TYPE = "v1.VolumeMount"
)

type VolumeMount struct {
	MountPath string `json:"mountPath,omitempty" yaml:"mount_path,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`
}

package model

const (
	CINDER_VOLUME_SOURCE_TYPE = "v1.CinderVolumeSource"
)

type CinderVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	VolumeID string `json:"volumeID,omitempty" yaml:"volume_id,omitempty"`
}

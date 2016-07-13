package model

const (
	EMPTY_DIR_VOLUME_SOURCE_TYPE = "v1.EmptyDirVolumeSource"
)

type EmptyDirVolumeSource struct {
	Medium string `json:"medium,omitempty" yaml:"medium,omitempty"`
}

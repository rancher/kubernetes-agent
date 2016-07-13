package model

const (
	HOST_PATH_VOLUME_SOURCE_TYPE = "v1.HostPathVolumeSource"
)

type HostPathVolumeSource struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

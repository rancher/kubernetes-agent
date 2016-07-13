package model

const (
	GLUSTERFS_VOLUME_SOURCE_TYPE = "v1.GlusterfsVolumeSource"
)

type GlusterfsVolumeSource struct {
	Endpoints string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`

	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`
}

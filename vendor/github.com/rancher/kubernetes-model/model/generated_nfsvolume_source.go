package model

const (
	NFSVOLUME_SOURCE_TYPE = "v1.NFSVolumeSource"
)

type NFSVolumeSource struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	Server string `json:"server,omitempty" yaml:"server,omitempty"`
}

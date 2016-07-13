package model

const (
	FCVOLUME_SOURCE_TYPE = "v1.FCVolumeSource"
)

type FCVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Lun int32 `json:"lun,omitempty" yaml:"lun,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	TargetWWNs []string `json:"targetWWNs,omitempty" yaml:"target_wwns,omitempty"`
}

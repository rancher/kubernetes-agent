package model

const (
	GCEPERSISTENT_DISK_VOLUME_SOURCE_TYPE = "v1.GCEPersistentDiskVolumeSource"
)

type GCEPersistentDiskVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Partition int32 `json:"partition,omitempty" yaml:"partition,omitempty"`

	PdName string `json:"pdName,omitempty" yaml:"pd_name,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`
}

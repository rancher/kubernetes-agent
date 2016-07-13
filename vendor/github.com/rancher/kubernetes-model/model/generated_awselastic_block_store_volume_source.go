package model

const (
	AWSELASTIC_BLOCK_STORE_VOLUME_SOURCE_TYPE = "v1.AWSElasticBlockStoreVolumeSource"
)

type AWSElasticBlockStoreVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Partition int32 `json:"partition,omitempty" yaml:"partition,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	VolumeID string `json:"volumeID,omitempty" yaml:"volume_id,omitempty"`
}

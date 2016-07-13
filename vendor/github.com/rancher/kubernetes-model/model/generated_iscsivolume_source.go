package model

const (
	ISCSIVOLUME_SOURCE_TYPE = "v1.ISCSIVolumeSource"
)

type ISCSIVolumeSource struct {
	FsType string `json:"fsType,omitempty" yaml:"fs_type,omitempty"`

	Iqn string `json:"iqn,omitempty" yaml:"iqn,omitempty"`

	IscsiInterface string `json:"iscsiInterface,omitempty" yaml:"iscsi_interface,omitempty"`

	Lun int32 `json:"lun,omitempty" yaml:"lun,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`

	TargetPortal string `json:"targetPortal,omitempty" yaml:"target_portal,omitempty"`
}

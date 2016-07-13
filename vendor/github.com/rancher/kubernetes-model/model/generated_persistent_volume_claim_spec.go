package model

const (
	PERSISTENT_VOLUME_CLAIM_SPEC_TYPE = "v1.PersistentVolumeClaimSpec"
)

type PersistentVolumeClaimSpec struct {
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" yaml:"access_modes,omitempty"`

	Resources *ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`

	VolumeName string `json:"volumeName,omitempty" yaml:"volume_name,omitempty"`
}

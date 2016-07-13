package model

const (
	PERSISTENT_VOLUME_CLAIM_STATUS_TYPE = "v1.PersistentVolumeClaimStatus"
)

type PersistentVolumeClaimStatus struct {
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" yaml:"access_modes,omitempty"`

	Capacity map[string]interface{} `json:"capacity,omitempty" yaml:"capacity,omitempty"`

	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`
}

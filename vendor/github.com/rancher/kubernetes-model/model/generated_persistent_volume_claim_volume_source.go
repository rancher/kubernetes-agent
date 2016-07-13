package model

const (
	PERSISTENT_VOLUME_CLAIM_VOLUME_SOURCE_TYPE = "v1.PersistentVolumeClaimVolumeSource"
)

type PersistentVolumeClaimVolumeSource struct {
	ClaimName string `json:"claimName,omitempty" yaml:"claim_name,omitempty"`

	ReadOnly bool `json:"readOnly,omitempty" yaml:"read_only,omitempty"`
}

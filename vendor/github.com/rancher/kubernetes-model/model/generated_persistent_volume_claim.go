package model

const (
	PERSISTENT_VOLUME_CLAIM_TYPE = "v1.PersistentVolumeClaim"
)

type PersistentVolumeClaim struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *PersistentVolumeClaimSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *PersistentVolumeClaimStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

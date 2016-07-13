package model

const (
	PERSISTENT_VOLUME_CLAIM_LIST_TYPE = "v1.PersistentVolumeClaimList"
)

type PersistentVolumeClaimList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []PersistentVolumeClaim `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

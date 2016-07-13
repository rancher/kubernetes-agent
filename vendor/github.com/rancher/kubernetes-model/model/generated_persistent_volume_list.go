package model

const (
	PERSISTENT_VOLUME_LIST_TYPE = "v1.PersistentVolumeList"
)

type PersistentVolumeList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []PersistentVolume `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

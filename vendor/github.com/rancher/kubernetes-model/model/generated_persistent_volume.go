package model

const (
	PERSISTENT_VOLUME_TYPE = "v1.PersistentVolume"
)

type PersistentVolume struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *PersistentVolumeSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *PersistentVolumeStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

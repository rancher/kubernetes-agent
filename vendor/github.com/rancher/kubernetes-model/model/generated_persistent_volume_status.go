package model

const (
	PERSISTENT_VOLUME_STATUS_TYPE = "v1.PersistentVolumeStatus"
)

type PersistentVolumeStatus struct {
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

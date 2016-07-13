package model

const (
	REPLICATION_CONTROLLER_TYPE = "v1.ReplicationController"
)

type ReplicationController struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *ReplicationControllerSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *ReplicationControllerStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

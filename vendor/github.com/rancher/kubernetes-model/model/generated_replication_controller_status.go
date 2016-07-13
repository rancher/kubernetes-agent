package model

const (
	REPLICATION_CONTROLLER_STATUS_TYPE = "v1.ReplicationControllerStatus"
)

type ReplicationControllerStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty" yaml:"observed_generation,omitempty"`

	Replicas int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
}

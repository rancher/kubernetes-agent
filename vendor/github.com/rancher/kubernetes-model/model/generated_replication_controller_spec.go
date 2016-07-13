package model

const (
	REPLICATION_CONTROLLER_SPEC_TYPE = "v1.ReplicationControllerSpec"
)

type ReplicationControllerSpec struct {
	Replicas int32 `json:"replicas,omitempty" yaml:"replicas,omitempty"`

	Selector map[string]interface{} `json:"selector,omitempty" yaml:"selector,omitempty"`

	Template *PodTemplateSpec `json:"template,omitempty" yaml:"template,omitempty"`
}

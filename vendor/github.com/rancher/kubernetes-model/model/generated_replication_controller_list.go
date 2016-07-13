package model

const (
	REPLICATION_CONTROLLER_LIST_TYPE = "v1.ReplicationControllerList"
)

type ReplicationControllerList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []ReplicationController `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

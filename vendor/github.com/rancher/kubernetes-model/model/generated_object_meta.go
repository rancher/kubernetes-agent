package model

const (
	OBJECT_META_TYPE = "v1.ObjectMeta"
)

type ObjectMeta struct {
	Annotations map[string]interface{} `json:"annotations,omitempty" yaml:"annotations,omitempty"`

	CreationTimestamp string `json:"creationTimestamp,omitempty" yaml:"creation_timestamp,omitempty"`

	DeletionGracePeriodSeconds int64 `json:"deletionGracePeriodSeconds,omitempty" yaml:"deletion_grace_period_seconds,omitempty"`

	DeletionTimestamp string `json:"deletionTimestamp,omitempty" yaml:"deletion_timestamp,omitempty"`

	GenerateName string `json:"generateName,omitempty" yaml:"generate_name,omitempty"`

	Generation int64 `json:"generation,omitempty" yaml:"generation,omitempty"`

	Labels map[string]interface{} `json:"labels,omitempty" yaml:"labels,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	ResourceVersion string `json:"resourceVersion,omitempty" yaml:"resource_version,omitempty"`

	SelfLink string `json:"selfLink,omitempty" yaml:"self_link,omitempty"`

	Uid string `json:"uid,omitempty" yaml:"uid,omitempty"`
}

package model

const (
	RESOURCE_QUOTA_TYPE = "v1.ResourceQuota"
)

type ResourceQuota struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Spec *ResourceQuotaSpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	Status *ResourceQuotaStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

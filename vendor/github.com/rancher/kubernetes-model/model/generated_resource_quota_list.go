package model

const (
	RESOURCE_QUOTA_LIST_TYPE = "v1.ResourceQuotaList"
)

type ResourceQuotaList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []ResourceQuota `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

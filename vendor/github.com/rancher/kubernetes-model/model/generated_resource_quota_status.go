package model

const (
	RESOURCE_QUOTA_STATUS_TYPE = "v1.ResourceQuotaStatus"
)

type ResourceQuotaStatus struct {
	Hard map[string]interface{} `json:"hard,omitempty" yaml:"hard,omitempty"`

	Used map[string]interface{} `json:"used,omitempty" yaml:"used,omitempty"`
}

package model

const (
	RESOURCE_QUOTA_SPEC_TYPE = "v1.ResourceQuotaSpec"
)

type ResourceQuotaSpec struct {
	Hard map[string]interface{} `json:"hard,omitempty" yaml:"hard,omitempty"`
}

package model

const (
	LIMIT_RANGE_ITEM_TYPE = "v1.LimitRangeItem"
)

type LimitRangeItem struct {
	Default map[string]interface{} `json:"default,omitempty" yaml:"default,omitempty"`

	DefaultRequest map[string]interface{} `json:"defaultRequest,omitempty" yaml:"default_request,omitempty"`

	Max map[string]interface{} `json:"max,omitempty" yaml:"max,omitempty"`

	MaxLimitRequestRatio map[string]interface{} `json:"maxLimitRequestRatio,omitempty" yaml:"max_limit_request_ratio,omitempty"`

	Min map[string]interface{} `json:"min,omitempty" yaml:"min,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

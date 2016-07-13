package model

const (
	RESOURCE_REQUIREMENTS_TYPE = "v1.ResourceRequirements"
)

type ResourceRequirements struct {
	Limits map[string]interface{} `json:"limits,omitempty" yaml:"limits,omitempty"`

	Requests map[string]interface{} `json:"requests,omitempty" yaml:"requests,omitempty"`
}

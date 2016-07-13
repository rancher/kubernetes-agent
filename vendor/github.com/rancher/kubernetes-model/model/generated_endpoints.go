package model

const (
	ENDPOINTS_TYPE = "v1.Endpoints"
)

type Endpoints struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Subsets []EndpointSubset `json:"subsets,omitempty" yaml:"subsets,omitempty"`
}

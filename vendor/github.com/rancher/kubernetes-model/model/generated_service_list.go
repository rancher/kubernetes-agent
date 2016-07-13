package model

const (
	SERVICE_LIST_TYPE = "v1.ServiceList"
)

type ServiceList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []Service `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

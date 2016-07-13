package model

const (
	NAMESPACE_LIST_TYPE = "v1.NamespaceList"
)

type NamespaceList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []Namespace `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

package model

const (
	POD_LIST_TYPE = "v1.PodList"
)

type PodList struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Items []Pod `json:"items,omitempty" yaml:"items,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ListMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

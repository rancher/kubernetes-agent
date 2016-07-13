package model

const (
	LIST_META_TYPE = "unversioned.ListMeta"
)

type ListMeta struct {
	ResourceVersion string `json:"resourceVersion,omitempty" yaml:"resource_version,omitempty"`

	SelfLink string `json:"selfLink,omitempty" yaml:"self_link,omitempty"`
}

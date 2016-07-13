package model

const (
	SECRET_TYPE = "v1.Secret"
)

type Secret struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	Data map[string]interface{} `json:"data,omitempty" yaml:"data,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

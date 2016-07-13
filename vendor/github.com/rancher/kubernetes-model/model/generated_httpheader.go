package model

const (
	HTTPHEADER_TYPE = "v1.HTTPHeader"
)

type HTTPHeader struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

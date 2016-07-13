package model

const (
	SECRET_KEY_SELECTOR_TYPE = "v1.SecretKeySelector"
)

type SecretKeySelector struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

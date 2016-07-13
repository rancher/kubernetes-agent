package model

const (
	CONFIG_MAP_KEY_SELECTOR_TYPE = "v1.ConfigMapKeySelector"
)

type ConfigMapKeySelector struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

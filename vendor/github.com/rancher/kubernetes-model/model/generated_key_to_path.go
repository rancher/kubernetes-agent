package model

const (
	KEY_TO_PATH_TYPE = "v1.KeyToPath"
)

type KeyToPath struct {
	Key string `json:"key,omitempty" yaml:"key,omitempty"`

	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

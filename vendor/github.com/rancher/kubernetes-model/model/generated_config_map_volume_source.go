package model

const (
	CONFIG_MAP_VOLUME_SOURCE_TYPE = "v1.ConfigMapVolumeSource"
)

type ConfigMapVolumeSource struct {
	Items []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

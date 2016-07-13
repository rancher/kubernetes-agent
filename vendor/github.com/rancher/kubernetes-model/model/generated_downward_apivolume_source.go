package model

const (
	DOWNWARD_APIVOLUME_SOURCE_TYPE = "v1.DownwardAPIVolumeSource"
)

type DownwardAPIVolumeSource struct {
	Items []DownwardAPIVolumeFile `json:"items,omitempty" yaml:"items,omitempty"`
}

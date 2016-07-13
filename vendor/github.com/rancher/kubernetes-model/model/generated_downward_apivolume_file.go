package model

const (
	DOWNWARD_APIVOLUME_FILE_TYPE = "v1.DownwardAPIVolumeFile"
)

type DownwardAPIVolumeFile struct {
	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" yaml:"field_ref,omitempty"`

	Path string `json:"path,omitempty" yaml:"path,omitempty"`
}

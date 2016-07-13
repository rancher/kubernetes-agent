package model

const (
	FLOCKER_VOLUME_SOURCE_TYPE = "v1.FlockerVolumeSource"
)

type FlockerVolumeSource struct {
	DatasetName string `json:"datasetName,omitempty" yaml:"dataset_name,omitempty"`
}

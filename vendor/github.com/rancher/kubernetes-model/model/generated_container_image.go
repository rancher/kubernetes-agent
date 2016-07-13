package model

const (
	CONTAINER_IMAGE_TYPE = "v1.ContainerImage"
)

type ContainerImage struct {
	Names []string `json:"names,omitempty" yaml:"names,omitempty"`

	SizeBytes int64 `json:"sizeBytes,omitempty" yaml:"size_bytes,omitempty"`
}

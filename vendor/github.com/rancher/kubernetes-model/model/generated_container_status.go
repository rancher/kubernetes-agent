package model

const (
	CONTAINER_STATUS_TYPE = "v1.ContainerStatus"
)

type ContainerStatus struct {
	ContainerID string `json:"containerID,omitempty" yaml:"container_id,omitempty"`

	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	ImageID string `json:"imageID,omitempty" yaml:"image_id,omitempty"`

	LastState *ContainerState `json:"lastState,omitempty" yaml:"last_state,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Ready bool `json:"ready,omitempty" yaml:"ready,omitempty"`

	RestartCount int32 `json:"restartCount,omitempty" yaml:"restart_count,omitempty"`

	State *ContainerState `json:"state,omitempty" yaml:"state,omitempty"`
}

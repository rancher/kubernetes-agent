package model

const (
	CONTAINER_STATE_TYPE = "v1.ContainerState"
)

type ContainerState struct {
	Running *ContainerStateRunning `json:"running,omitempty" yaml:"running,omitempty"`

	Terminated *ContainerStateTerminated `json:"terminated,omitempty" yaml:"terminated,omitempty"`

	Waiting *ContainerStateWaiting `json:"waiting,omitempty" yaml:"waiting,omitempty"`
}

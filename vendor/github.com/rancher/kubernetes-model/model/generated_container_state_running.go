package model

const (
	CONTAINER_STATE_RUNNING_TYPE = "v1.ContainerStateRunning"
)

type ContainerStateRunning struct {
	StartedAt string `json:"startedAt,omitempty" yaml:"started_at,omitempty"`
}

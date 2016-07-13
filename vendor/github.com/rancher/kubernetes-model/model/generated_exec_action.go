package model

const (
	EXEC_ACTION_TYPE = "v1.ExecAction"
)

type ExecAction struct {
	Command []string `json:"command,omitempty" yaml:"command,omitempty"`
}

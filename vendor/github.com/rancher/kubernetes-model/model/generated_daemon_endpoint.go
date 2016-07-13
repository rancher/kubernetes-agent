package model

const (
	DAEMON_ENDPOINT_TYPE = "v1.DaemonEndpoint"
)

type DaemonEndpoint struct {
	Port int32 `json:"port,omitempty" yaml:"port,omitempty"`
}

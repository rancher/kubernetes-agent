package model

const (
	TCPSOCKET_ACTION_TYPE = "v1.TCPSocketAction"
)

type TCPSocketAction struct {
	Port interface{} `json:"port,omitempty" yaml:"port,omitempty"`
}

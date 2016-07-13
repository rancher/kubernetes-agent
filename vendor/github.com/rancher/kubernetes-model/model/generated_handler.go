package model

const (
	HANDLER_TYPE = "v1.Handler"
)

type Handler struct {
	Exec *ExecAction `json:"exec,omitempty" yaml:"exec,omitempty"`

	HttpGet *HTTPGetAction `json:"httpGet,omitempty" yaml:"http_get,omitempty"`

	TcpSocket *TCPSocketAction `json:"tcpSocket,omitempty" yaml:"tcp_socket,omitempty"`
}

package model

const (
	PROBE_TYPE = "v1.Probe"
)

type Probe struct {
	Exec *ExecAction `json:"exec,omitempty" yaml:"exec,omitempty"`

	FailureThreshold int32 `json:"failureThreshold,omitempty" yaml:"failure_threshold,omitempty"`

	HttpGet *HTTPGetAction `json:"httpGet,omitempty" yaml:"http_get,omitempty"`

	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" yaml:"initial_delay_seconds,omitempty"`

	PeriodSeconds int32 `json:"periodSeconds,omitempty" yaml:"period_seconds,omitempty"`

	SuccessThreshold int32 `json:"successThreshold,omitempty" yaml:"success_threshold,omitempty"`

	TcpSocket *TCPSocketAction `json:"tcpSocket,omitempty" yaml:"tcp_socket,omitempty"`

	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" yaml:"timeout_seconds,omitempty"`
}

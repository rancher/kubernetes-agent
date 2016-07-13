package model

const (
	SERVICE_PORT_TYPE = "v1.ServicePort"
)

type ServicePort struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	NodePort int32 `json:"nodePort,omitempty" yaml:"node_port,omitempty"`

	Port int32 `json:"port,omitempty" yaml:"port,omitempty"`

	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`

	TargetPort interface{} `json:"targetPort,omitempty" yaml:"target_port,omitempty"`
}

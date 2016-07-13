package model

const (
	CONTAINER_PORT_TYPE = "v1.ContainerPort"
)

type ContainerPort struct {
	ContainerPort int32 `json:"containerPort,omitempty" yaml:"container_port,omitempty"`

	HostIP string `json:"hostIP,omitempty" yaml:"host_ip,omitempty"`

	HostPort int32 `json:"hostPort,omitempty" yaml:"host_port,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

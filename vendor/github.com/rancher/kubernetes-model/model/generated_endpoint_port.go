package model

const (
	ENDPOINT_PORT_TYPE = "v1.EndpointPort"
)

type EndpointPort struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Port int32 `json:"port,omitempty" yaml:"port,omitempty"`

	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

package model

const (
	LOAD_BALANCER_INGRESS_TYPE = "v1.LoadBalancerIngress"
)

type LoadBalancerIngress struct {
	Hostname string `json:"hostname,omitempty" yaml:"hostname,omitempty"`

	Ip string `json:"ip,omitempty" yaml:"ip,omitempty"`
}

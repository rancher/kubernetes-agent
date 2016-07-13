package model

const (
	LOAD_BALANCER_STATUS_TYPE = "v1.LoadBalancerStatus"
)

type LoadBalancerStatus struct {
	Ingress []LoadBalancerIngress `json:"ingress,omitempty" yaml:"ingress,omitempty"`
}

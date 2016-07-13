package model

const (
	SERVICE_STATUS_TYPE = "v1.ServiceStatus"
)

type ServiceStatus struct {
	LoadBalancer *LoadBalancerStatus `json:"loadBalancer,omitempty" yaml:"load_balancer,omitempty"`
}

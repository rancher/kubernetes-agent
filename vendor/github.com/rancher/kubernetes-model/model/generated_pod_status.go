package model

const (
	POD_STATUS_TYPE = "v1.PodStatus"
)

type PodStatus struct {
	Conditions []PodCondition `json:"conditions,omitempty" yaml:"conditions,omitempty"`

	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty" yaml:"container_statuses,omitempty"`

	HostIP string `json:"hostIP,omitempty" yaml:"host_ip,omitempty"`

	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	Phase string `json:"phase,omitempty" yaml:"phase,omitempty"`

	PodIP string `json:"podIP,omitempty" yaml:"pod_ip,omitempty"`

	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`

	StartTime string `json:"startTime,omitempty" yaml:"start_time,omitempty"`
}

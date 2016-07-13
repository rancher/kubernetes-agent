package model

const (
	POD_SPEC_TYPE = "v1.PodSpec"
)

type PodSpec struct {
	ActiveDeadlineSeconds int64 `json:"activeDeadlineSeconds,omitempty" yaml:"active_deadline_seconds,omitempty"`

	Containers []Container `json:"containers,omitempty" yaml:"containers,omitempty"`

	DnsPolicy string `json:"dnsPolicy,omitempty" yaml:"dns_policy,omitempty"`

	HostIPC bool `json:"hostIPC,omitempty" yaml:"host_ipc,omitempty"`

	HostNetwork bool `json:"hostNetwork,omitempty" yaml:"host_network,omitempty"`

	HostPID bool `json:"hostPID,omitempty" yaml:"host_pid,omitempty"`

	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty" yaml:"image_pull_secrets,omitempty"`

	NodeName string `json:"nodeName,omitempty" yaml:"node_name,omitempty"`

	NodeSelector map[string]interface{} `json:"nodeSelector,omitempty" yaml:"node_selector,omitempty"`

	RestartPolicy string `json:"restartPolicy,omitempty" yaml:"restart_policy,omitempty"`

	SecurityContext *PodSecurityContext `json:"securityContext,omitempty" yaml:"security_context,omitempty"`

	ServiceAccount string `json:"serviceAccount,omitempty" yaml:"service_account,omitempty"`

	ServiceAccountName string `json:"serviceAccountName,omitempty" yaml:"service_account_name,omitempty"`

	TerminationGracePeriodSeconds int64 `json:"terminationGracePeriodSeconds,omitempty" yaml:"termination_grace_period_seconds,omitempty"`

	Volumes []Volume `json:"volumes,omitempty" yaml:"volumes,omitempty"`
}

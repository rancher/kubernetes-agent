package model

const (
	NODE_SYSTEM_INFO_TYPE = "v1.NodeSystemInfo"
)

type NodeSystemInfo struct {
	BootID string `json:"bootID,omitempty" yaml:"boot_id,omitempty"`

	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty" yaml:"container_runtime_version,omitempty"`

	KernelVersion string `json:"kernelVersion,omitempty" yaml:"kernel_version,omitempty"`

	KubeProxyVersion string `json:"kubeProxyVersion,omitempty" yaml:"kube_proxy_version,omitempty"`

	KubeletVersion string `json:"kubeletVersion,omitempty" yaml:"kubelet_version,omitempty"`

	MachineID string `json:"machineID,omitempty" yaml:"machine_id,omitempty"`

	OsImage string `json:"osImage,omitempty" yaml:"os_image,omitempty"`

	SystemUUID string `json:"systemUUID,omitempty" yaml:"system_uuid,omitempty"`
}

package model

const (
	CONTAINER_TYPE = "v1.Container"
)

type Container struct {
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`

	Command []string `json:"command,omitempty" yaml:"command,omitempty"`

	Env []EnvVar `json:"env,omitempty" yaml:"env,omitempty"`

	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	ImagePullPolicy string `json:"imagePullPolicy,omitempty" yaml:"image_pull_policy,omitempty"`

	Lifecycle *Lifecycle `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`

	LivenessProbe *Probe `json:"livenessProbe,omitempty" yaml:"liveness_probe,omitempty"`

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Ports []ContainerPort `json:"ports,omitempty" yaml:"ports,omitempty"`

	ReadinessProbe *Probe `json:"readinessProbe,omitempty" yaml:"readiness_probe,omitempty"`

	Resources *ResourceRequirements `json:"resources,omitempty" yaml:"resources,omitempty"`

	SecurityContext *SecurityContext `json:"securityContext,omitempty" yaml:"security_context,omitempty"`

	Stdin bool `json:"stdin,omitempty" yaml:"stdin,omitempty"`

	StdinOnce bool `json:"stdinOnce,omitempty" yaml:"stdin_once,omitempty"`

	TerminationMessagePath string `json:"terminationMessagePath,omitempty" yaml:"termination_message_path,omitempty"`

	Tty bool `json:"tty,omitempty" yaml:"tty,omitempty"`

	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" yaml:"volume_mounts,omitempty"`

	WorkingDir string `json:"workingDir,omitempty" yaml:"working_dir,omitempty"`
}

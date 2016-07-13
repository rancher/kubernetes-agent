package model

const (
	SECURITY_CONTEXT_TYPE = "v1.SecurityContext"
)

type SecurityContext struct {
	Capabilities *Capabilities `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`

	Privileged bool `json:"privileged,omitempty" yaml:"privileged,omitempty"`

	ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem,omitempty" yaml:"read_only_root_filesystem,omitempty"`

	RunAsNonRoot bool `json:"runAsNonRoot,omitempty" yaml:"run_as_non_root,omitempty"`

	RunAsUser int64 `json:"runAsUser,omitempty" yaml:"run_as_user,omitempty"`

	SeLinuxOptions *SELinuxOptions `json:"seLinuxOptions,omitempty" yaml:"se_linux_options,omitempty"`
}

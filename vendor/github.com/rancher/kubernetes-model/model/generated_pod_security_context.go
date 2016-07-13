package model

const (
	POD_SECURITY_CONTEXT_TYPE = "v1.PodSecurityContext"
)

type PodSecurityContext struct {
	FsGroup int64 `json:"fsGroup,omitempty" yaml:"fs_group,omitempty"`

	RunAsNonRoot bool `json:"runAsNonRoot,omitempty" yaml:"run_as_non_root,omitempty"`

	RunAsUser int64 `json:"runAsUser,omitempty" yaml:"run_as_user,omitempty"`

	SeLinuxOptions *SELinuxOptions `json:"seLinuxOptions,omitempty" yaml:"se_linux_options,omitempty"`

	SupplementalGroups []int64 `json:"supplementalGroups,omitempty" yaml:"supplemental_groups,omitempty"`
}

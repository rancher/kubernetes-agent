package model

const (
	SELINUX_OPTIONS_TYPE = "v1.SELinuxOptions"
)

type SELinuxOptions struct {
	Level string `json:"level,omitempty" yaml:"level,omitempty"`

	Role string `json:"role,omitempty" yaml:"role,omitempty"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	User string `json:"user,omitempty" yaml:"user,omitempty"`
}

package model

const (
	SERVICE_ACCOUNT_TYPE = "v1.ServiceAccount"
)

type ServiceAccount struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"api_version,omitempty"`

	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty" yaml:"image_pull_secrets,omitempty"`

	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	Metadata *ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	Secrets []ObjectReference `json:"secrets,omitempty" yaml:"secrets,omitempty"`
}

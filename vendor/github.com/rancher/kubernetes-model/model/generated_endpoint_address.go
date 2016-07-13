package model

const (
	ENDPOINT_ADDRESS_TYPE = "v1.EndpointAddress"
)

type EndpointAddress struct {
	Ip string `json:"ip,omitempty" yaml:"ip,omitempty"`

	TargetRef *ObjectReference `json:"targetRef,omitempty" yaml:"target_ref,omitempty"`
}

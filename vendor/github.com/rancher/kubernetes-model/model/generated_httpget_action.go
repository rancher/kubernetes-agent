package model

const (
	HTTPGET_ACTION_TYPE = "v1.HTTPGetAction"
)

type HTTPGetAction struct {
	Host string `json:"host,omitempty" yaml:"host,omitempty"`

	HttpHeaders []HTTPHeader `json:"httpHeaders,omitempty" yaml:"http_headers,omitempty"`

	Path string `json:"path,omitempty" yaml:"path,omitempty"`

	Port interface{} `json:"port,omitempty" yaml:"port,omitempty"`

	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
}

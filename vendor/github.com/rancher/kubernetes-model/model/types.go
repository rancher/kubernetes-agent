package model

type Schema struct {
	APIVerion string           `json:"apiVersion,omitempty"`
	Models    map[string]Model `json:"models,omitempty"`
}

type Model struct {
	ID         string              `json:"id,omitempty"`
	Required   []string            `json:"required,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Items       Items  `json:"items,omitempty"`
	Ref         string `json:"$ref,omitempty"`
	Format      string `json:"format,omitempty"`
}

type Items struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

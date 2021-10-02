package goai

// Tags is specified by OpenAPI/Swagger 3.0 standard.
type Tags []Tag

// Tag is specified by OpenAPI/Swagger 3.0 standard.
type Tag struct {
	Name         string        `json:"name,omitempty"         yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty"  yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

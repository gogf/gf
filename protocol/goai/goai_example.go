package goai

// Example is specified by OpenAPI/Swagger 3.0 standard.
type Example struct {
	Summary       string      `json:"summary,omitempty"       yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty"   yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty"         yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

type Examples map[string]*ExampleRef

type ExampleRef struct {
	Ref   string
	Value *Example
}

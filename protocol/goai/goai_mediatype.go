package goai

// MediaType is specified by OpenAPI/Swagger 3.0 standard.
type MediaType struct {
	Schema   *SchemaRef           `json:"schema,omitempty"   yaml:"schema,omitempty"`
	Example  interface{}          `json:"example,omitempty"  yaml:"example,omitempty"`
	Examples Examples             `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// Content is specified by OpenAPI/Swagger 3.0 standard.
type Content map[string]MediaType

// Encoding is specified by OpenAPI/Swagger 3.0 standard.
type Encoding struct {
	ContentType   string  `json:"contentType,omitempty"   yaml:"contentType,omitempty"`
	Headers       Headers `json:"headers,omitempty"       yaml:"headers,omitempty"`
	Style         string  `json:"style,omitempty"         yaml:"style,omitempty"`
	Explode       *bool   `json:"explode,omitempty"       yaml:"explode,omitempty"`
	AllowReserved bool    `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

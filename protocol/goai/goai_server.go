package goai

// Server is specified by OpenAPI/Swagger standard version 3.0.
type Server struct {
	URL         string                     `json:"url"                   yaml:"url"`
	Description string                     `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"   yaml:"variables,omitempty"`
}

// ServerVariable is specified by OpenAPI/Swagger standard version 3.0.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"        yaml:"enum,omitempty"`
	Default     string   `json:"default,omitempty"     yaml:"default,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// Servers is specified by OpenAPI/Swagger standard version 3.0.
type Servers []Server

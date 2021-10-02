package goai

// Info is specified by OpenAPI/Swagger standard version 3.0.
type Info struct {
	Title          string   `json:"title"                    yaml:"title"`
	Description    string   `json:"description,omitempty"    yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"        yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty"        yaml:"license,omitempty"`
	Version        string   `json:"version"                  yaml:"version"`
}

// Contact is specified by OpenAPI/Swagger standard version 3.0.
type Contact struct {
	Name  string `json:"name,omitempty"  yaml:"name,omitempty"`
	URL   string `json:"url,omitempty"   yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License is specified by OpenAPI/Swagger standard version 3.0.
type License struct {
	Name string `json:"name"          yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

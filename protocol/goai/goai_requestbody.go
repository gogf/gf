package goai

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
)

// RequestBody is specified by OpenAPI/Swagger 3.0 standard.
type RequestBody struct {
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool    `json:"required,omitempty"    yaml:"required,omitempty"`
	Content     Content `json:"content,omitempty"     yaml:"content,omitempty"`
}

type RequestBodyRef struct {
	Ref   string
	Value *RequestBody
}

func (r RequestBodyRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return []byte(fmt.Sprintf(`{"$ref":"#/components/schemas/%s"}`, r.Ref)), nil
	}
	return json.Marshal(r.Value)
}

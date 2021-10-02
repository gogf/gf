// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/internal/json"
)

// Response is specified by OpenAPI/Swagger 3.0 standard.
type Response struct {
	Description string  `json:"description"           yaml:"description"`
	Headers     Headers `json:"headers,omitempty"     yaml:"headers,omitempty"`
	Content     Content `json:"content,omitempty"     yaml:"content,omitempty"`
	Links       Links   `json:"links,omitempty"       yaml:"links,omitempty"`
}

// Responses is specified by OpenAPI/Swagger 3.0 standard.
type Responses map[string]ResponseRef

type ResponseRef struct {
	Ref   string
	Value *Response
}

func (r ResponseRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

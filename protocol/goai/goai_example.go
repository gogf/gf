// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/internal/json"
)

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

func (r ExampleRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

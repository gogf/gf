// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

// Components is specified by OpenAPI/Swagger standard version 3.0.
type Components struct {
	Schemas         Schemas         `json:"schemas,omitempty"         yaml:"schemas,omitempty"`
	Parameters      ParametersMap   `json:"parameters,omitempty"      yaml:"parameters,omitempty"`
	Headers         Headers         `json:"headers,omitempty"         yaml:"headers,omitempty"`
	RequestBodies   RequestBodies   `json:"requestBodies,omitempty"   yaml:"requestBodies,omitempty"`
	Responses       Responses       `json:"responses,omitempty"       yaml:"responses,omitempty"`
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	Examples        Examples        `json:"examples,omitempty"        yaml:"examples,omitempty"`
	Links           Links           `json:"links,omitempty"           yaml:"links,omitempty"`
	Callbacks       Callbacks       `json:"callbacks,omitempty"       yaml:"callbacks,omitempty"`
}

type ParametersMap map[string]*ParameterRef

type RequestBodies map[string]*RequestBodyRef

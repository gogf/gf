// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

// Tags is specified by OpenAPI/Swagger 3.0 standard.
type Tags []Tag

// Tag is specified by OpenAPI/Swagger 3.0 standard.
type Tag struct {
	Name         string        `json:"name,omitempty"         yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty"  yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

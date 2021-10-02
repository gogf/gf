// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

// Operation represents "operation" specified by OpenAPI/Swagger 3.0 standard.
type Operation struct {
	Tags         []string              `json:"tags,omitempty"         yaml:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"      yaml:"summary,omitempty"`
	Description  string                `json:"description,omitempty"  yaml:"description,omitempty"`
	OperationID  string                `json:"operationId,omitempty"  yaml:"operationId,omitempty"`
	Parameters   Parameters            `json:"parameters,omitempty"   yaml:"parameters,omitempty"`
	RequestBody  RequestBodyRef        `json:"requestBody,omitempty"  yaml:"requestBody,omitempty"`
	Responses    Responses             `json:"responses"              yaml:"responses"`
	Deprecated   bool                  `json:"deprecated,omitempty"   yaml:"deprecated,omitempty"`
	Callbacks    *Callbacks            `json:"callbacks,omitempty"    yaml:"callbacks,omitempty"`
	Security     *SecurityRequirements `json:"security,omitempty"     yaml:"security,omitempty"`
	Servers      *Servers              `json:"servers,omitempty"      yaml:"servers,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/internal/json"
)

// Link is specified by OpenAPI/Swagger standard version 3.0.
type Link struct {
	OperationID  string         `json:"operationId,omitempty"`
	OperationRef string         `json:"operationRef,omitempty"`
	Description  string         `json:"description,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	Server       *Server        `json:"server,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
}

type Links map[string]LinkRef

type LinkRef struct {
	Ref   string
	Value *Link
}

func (r LinkRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

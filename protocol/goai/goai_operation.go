// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

// Operation represents "operation" specified by OpenAPI/Swagger 3.0 standard.
type Operation struct {
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   Parameters            `json:"parameters,omitempty"`
	RequestBody  *RequestBodyRef       `json:"requestBody,omitempty"`
	Responses    Responses             `json:"responses"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Callbacks    *Callbacks            `json:"callbacks,omitempty"`
	Security     *SecurityRequirements `json:"security,omitempty"`
	Servers      *Servers              `json:"servers,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
	XExtensions  XExtensions           `json:"-"`
}

func (oai *OpenApiV3) tagMapToOperation(tagMap map[string]string, operation *Operation) error {
	var mergedTagMap = oai.fileMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, operation); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Operation failed`)
	}
	oai.tagMapToXExtensions(mergedTagMap, operation.XExtensions)
	return nil
}

func (o Operation) MarshalJSON() ([]byte, error) {
	var (
		b   []byte
		m   map[string]json.RawMessage
		err error
	)
	type tempOperation Operation // To prevent JSON marshal recursion error.
	if b, err = json.Marshal(tempOperation(o)); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range o.XExtensions {
		if b, err = json.Marshal(v); err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

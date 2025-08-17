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

// Parameter is specified by OpenAPI/Swagger 3.0 standard.
// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md#parameterObject
type Parameter struct {
	Name            string      `json:"name,omitempty"`
	In              string      `json:"in,omitempty"`
	Description     string      `json:"description,omitempty"`
	Style           string      `json:"style,omitempty"`
	Explode         *bool       `json:"explode,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"`
	Required        bool        `json:"required,omitempty"`
	Schema          *SchemaRef  `json:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty"`
	Examples        *Examples   `json:"examples,omitempty"`
	Content         *Content    `json:"content,omitempty"`
	XExtensions     XExtensions `json:"-"`
}

func (oai *OpenApiV3) tagMapToParameter(tagMap map[string]string, parameter *Parameter) error {
	var mergedTagMap = oai.fillMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, parameter); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Parameter failed`)
	}
	oai.tagMapToXExtensions(mergedTagMap, parameter.XExtensions)
	return nil
}

func (p Parameter) MarshalJSON() ([]byte, error) {
	var (
		b   []byte
		m   map[string]json.RawMessage
		err error
	)
	type tempParameter Parameter // To prevent JSON marshal recursion error.
	if b, err = json.Marshal(tempParameter(p)); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range p.XExtensions {
		if b, err = json.Marshal(v); err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

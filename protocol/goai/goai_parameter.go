// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
)

// Parameter is specified by OpenAPI/Swagger 3.0 standard.
// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md#parameterObject
type Parameter struct {
	Name            string      `json:"name,omitempty"            yaml:"name,omitempty"`
	In              string      `json:"in,omitempty"              yaml:"in,omitempty"`
	Description     string      `json:"description,omitempty"     yaml:"description,omitempty"`
	Style           string      `json:"style,omitempty"           yaml:"style,omitempty"`
	Explode         *bool       `json:"explode,omitempty"         yaml:"explode,omitempty"`
	AllowEmptyValue bool        `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	AllowReserved   bool        `json:"allowReserved,omitempty"   yaml:"allowReserved,omitempty"`
	Deprecated      bool        `json:"deprecated,omitempty"      yaml:"deprecated,omitempty"`
	Required        bool        `json:"required,omitempty"        yaml:"required,omitempty"`
	Schema          *SchemaRef  `json:"schema,omitempty"          yaml:"schema,omitempty"`
	Example         interface{} `json:"example,omitempty"         yaml:"example,omitempty"`
	Examples        *Examples   `json:"examples,omitempty"        yaml:"examples,omitempty"`
	Content         *Content    `json:"content,omitempty"         yaml:"content,omitempty"`
}

// Parameters is specified by OpenAPI/Swagger 3.0 standard.
type Parameters []ParameterRef

type ParameterRef struct {
	Ref   string
	Value *Parameter
}

func (oai *OpenApiV3) newParameterRefWithStructMethod(field *structs.Field) (*ParameterRef, error) {
	var (
		inTagValue = field.Tag(TagNameIn)
	)
	if inTagValue == "" {
		return nil, nil
	}
	var (
		tagMap    = field.TagMap()
		parameter = &Parameter{
			Name: field.TagJsonName(),
		}
	)
	if parameter.Name == "" {
		parameter.Name = field.Name()
	}
	if len(tagMap) > 0 {
		err := gconv.Struct(tagMap, parameter)
		if err != nil {
			return nil, gerror.WrapCodef(gcode.CodeInternalError, err, `mapping struct tags to Parameter failed`)
		}
	}
	// Required for path parameter.
	switch inTagValue {
	case ParameterInPath:
		parameter.Required = true
	case ParameterInCookie, ParameterInHeader, ParameterInQuery:
		// Nothing to do.

	default:
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid value "%s" for In`, inTagValue)
	}
	// Necessary schema or content.
	schemaRef, err := oai.newSchemaRefWithGolangType(field.Type(), tagMap)
	if err != nil {
		return nil, err
	}
	parameter.Schema = schemaRef

	return &ParameterRef{
		Ref:   "",
		Value: parameter,
	}, nil
}

func (r ParameterRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
)

// Parameters is specified by OpenAPI/Swagger 3.0 standard.
type Parameters []ParameterRef

type ParameterRef struct {
	Ref   string
	Value *Parameter
}

func (oai *OpenApiV3) newParameterRefWithStructMethod(field gstructs.Field, path, method string) (*ParameterRef, error) {
	var (
		tagMap    = field.TagMap()
		parameter = &Parameter{
			Name:        field.TagJsonName(),
			XExtensions: make(XExtensions),
		}
	)
	if parameter.Name == "" {
		parameter.Name = field.Name()
	}
	if len(tagMap) > 0 {
		if err := oai.tagMapToParameter(tagMap, parameter); err != nil {
			return nil, err
		}
	}
	if parameter.In == "" {
		// Automatically detect its "in" attribute.
		if gstr.ContainsI(path, fmt.Sprintf(`{%s}`, parameter.Name)) {
			parameter.In = ParameterInPath
		} else {
			// Default the parameter input to "query" if method is "GET/DELETE".
			switch gstr.ToUpper(method) {
			case HttpMethodGet, HttpMethodDelete:
				parameter.In = ParameterInQuery

			default:
				return nil, nil
			}
		}
	}

	switch parameter.In {
	case ParameterInPath:
		// Required for path parameter.
		parameter.Required = true

	case ParameterInCookie, ParameterInHeader, ParameterInQuery:

	default:
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid tag value "%s" for In`, parameter.In)
	}
	// Necessary schema or content.
	schemaRef, err := oai.newSchemaRefWithGolangType(field.Type().Type, tagMap)
	if err != nil {
		return nil, err
	}
	parameter.Schema = schemaRef

	// Required check.
	if parameter.Schema.Value != nil && parameter.Schema.Value.Pattern != "" {
		if gset.NewStrSetFrom(gstr.Split(parameter.Schema.Value.Pattern, "|")).Contains(patternKeyForRequired) {
			parameter.Required = true
		}
	}

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

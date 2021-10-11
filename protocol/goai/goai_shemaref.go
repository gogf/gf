// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
	"reflect"
)

type SchemaRefs []SchemaRef

type SchemaRef struct {
	Ref   string
	Value *Schema
}

func (oai *OpenApiV3) newSchemaRefWithGolangType(golangType reflect.Type, tagMap map[string]string) (*SchemaRef, error) {
	var (
		oaiType   = oai.golangTypeToOAIType(golangType)
		oaiFormat = oai.golangTypeToOAIFormat(golangType)
		schemaRef = &SchemaRef{}
		schema    = &Schema{
			Type:   oaiType,
			Format: oaiFormat,
		}
	)
	if len(tagMap) > 0 {
		if err := gconv.Struct(tagMap, schema); err != nil {
			return nil, gerror.WrapCode(gcode.CodeInternalError, err, `mapping struct tags to Schema failed`)
		}
	}
	schemaRef.Value = schema
	switch oaiType {
	case
		TypeNumber,
		TypeString,
		TypeBoolean:
		// Nothing to do.

	case
		TypeArray:
		subSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
		if err != nil {
			return nil, err
		}
		schema.Items = subSchemaRef

	case
		TypeObject:
		var (
			structTypeName = golangTypeToSchemaName(golangType)
		)
		// Specially for map type.
		if golangType.Kind() == reflect.Map {
			subSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
			if err != nil {
				return nil, err
			}
			schema.AdditionalProperties = subSchemaRef
			return schemaRef, nil
		}
		// Normal struct object.
		if _, ok := oai.Components.Schemas[structTypeName]; !ok {
			if err := oai.addSchema(reflect.New(golangType).Interface()); err != nil {
				return nil, err
			}
		}
		schemaRef.Ref = structTypeName
		schemaRef.Value = nil
	}
	return schemaRef, nil
}

func (r SchemaRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

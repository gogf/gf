// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
)

type SchemaRefs []SchemaRef

type SchemaRef struct {
	Ref   string
	Value *Schema
}

// isEmbeddedStructDefine checks and returns whether given golang type is embedded struct definition, like:
//
//	struct A struct{
//	    B struct{
//	        // ...
//	    }
//	}
//
// The `B` in `A` is called `embedded struct definition`.
func (oai *OpenApiV3) isEmbeddedStructDefinition(golangType reflect.Type) bool {
	s := golangType.String()
	return gstr.Contains(s, `struct {`)
}

// newSchemaRefWithGolangType creates a new Schema and returns its SchemaRef.
func (oai *OpenApiV3) newSchemaRefWithGolangType(golangType reflect.Type, tagMap map[string]string) (*SchemaRef, error) {
	var (
		err       error
		oaiType   = oai.golangTypeToOAIType(golangType)
		oaiFormat = oai.golangTypeToOAIFormat(golangType)
		typeName  = golangType.Name()
		pkgPath   = golangType.PkgPath()
		schemaRef = &SchemaRef{}
		schema    = &Schema{
			Type:        oaiType,
			Format:      oaiFormat,
			XExtensions: make(XExtensions),
		}
	)
	if pkgPath == "" && golangType.Kind() == reflect.Ptr {
		pkgPath = golangType.Elem().PkgPath()
		typeName = golangType.Elem().Name()
	}

	// Type enums.
	var typeId = fmt.Sprintf(`%s.%s`, pkgPath, typeName)
	if enums := gtag.GetEnumsByType(typeId); enums != "" {
		schema.Enum = make([]interface{}, 0)
		if err = json.Unmarshal([]byte(enums), &schema.Enum); err != nil {
			return nil, err
		}
	}

	if len(tagMap) > 0 {
		if err := oai.tagMapToSchema(tagMap, schema); err != nil {
			return nil, err
		}
	}
	schemaRef.Value = schema
	switch oaiType {
	case TypeString:
	// Nothing to do.
	case TypeInteger:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Int64(schemaRef.Value.Default)
		}
		// keep the default value as nil.
	case TypeNumber:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Float64(schemaRef.Value.Default)
		}
		// keep the default value as nil.
	case TypeBoolean:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Bool(schemaRef.Value.Default)
		}
		// keep the default value as nil.
	case
		TypeArray:
		subSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
		if err != nil {
			return nil, err
		}
		schema.Items = subSchemaRef
		if len(schema.Enum) > 0 {
			schema.Items.Value.Enum = schema.Enum
			schema.Enum = nil
		}

	case
		TypeObject:
		for golangType.Kind() == reflect.Ptr {
			golangType = golangType.Elem()
		}
		switch golangType.Kind() {
		case reflect.Map:
			// Specially for map type.
			subSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
			if err != nil {
				return nil, err
			}
			schema.AdditionalProperties = subSchemaRef
			return schemaRef, nil

		case reflect.Interface:
			// Specially for interface type.
			var (
				structTypeName = oai.golangTypeToSchemaName(golangType)
			)
			if oai.Components.Schemas.Get(structTypeName) == nil {
				if err := oai.addSchema(reflect.New(golangType).Interface()); err != nil {
					return nil, err
				}
			}
			schemaRef.Ref = structTypeName
			schemaRef.Value = nil

		default:
			golangTypeInstance := reflect.New(golangType).Elem().Interface()
			if oai.isEmbeddedStructDefinition(golangType) {
				schema, err = oai.structToSchema(golangTypeInstance)
				if err != nil {
					return nil, err
				}
				schemaRef.Ref = ""
				schemaRef.Value = schema
			} else {
				var structTypeName = oai.golangTypeToSchemaName(golangType)
				if oai.Components.Schemas.Get(structTypeName) == nil {
					if err := oai.addSchema(golangTypeInstance); err != nil {
						return nil, err
					}
				}
				schemaRef.Ref = structTypeName
				schemaRef.Value = nil
			}
		}
	}
	return schemaRef, nil
}

func (r SchemaRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

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
	Ref         string
	Description string
	Value       *Schema
}

// isEmbeddedStructDefinition checks and returns whether given golang type is embedded struct definition, like:
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
	if pkgPath == "" {
		switch golangType.Kind() {
		case reflect.Pointer, reflect.Array, reflect.Slice:
			pkgPath = golangType.Elem().PkgPath()
			typeName = golangType.Elem().Name()
		default:
		}
	}

	// Type enums.
	var typeId = fmt.Sprintf(`%s.%s`, pkgPath, typeName)
	enumItems, err := gtag.GetEnumItemsByType(typeId)
	if err != nil {
		return nil, err
	}
	schema.Enum = make([]any, 0, len(enumItems))
	for _, enumItem := range enumItems {
		schema.Enum = append(schema.Enum, enumItem.Value)
	}

	if len(tagMap) > 0 {
		if err = oai.tagMapToSchema(tagMap, schema); err != nil {
			return nil, err
		}
		if oaiType == TypeArray && schema.Type == TypeFile {
			schema.Type = TypeArray
		}
	}
	schemaRef.Value = schema
	switch schema.Type {
	case TypeString, TypeFile:
	// Nothing to do.
	case TypeInteger:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Int64(schemaRef.Value.Default)
		}
		// keep the default value as nil.

		// example value needs to be converted just like default value
		if schemaRef.Value.Example != nil {
			schemaRef.Value.Example = gconv.Int64(schemaRef.Value.Example)
		}
		// keep the example value as nil.
	case TypeNumber:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Float64(schemaRef.Value.Default)
		}
		// keep the default value as nil.

		// example value needs to be converted just like default value
		if schemaRef.Value.Example != nil {
			schemaRef.Value.Example = gconv.Float64(schemaRef.Value.Example)
		}
		// keep the example value as nil.
	case TypeBoolean:
		if schemaRef.Value.Default != nil {
			schemaRef.Value.Default = gconv.Bool(schemaRef.Value.Default)
		}
		// keep the default value as nil.

		// example value needs to be converted just like default value
		if schemaRef.Value.Example != nil {
			schemaRef.Value.Example = gconv.Bool(schemaRef.Value.Example)
		}
		// keep the example value as nil.
	case TypeArray:
		subSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
		if err != nil {
			return nil, err
		}
		schema.Items = subSchemaRef
		if len(schema.Enum) > 0 {
			schema.Items.Value.Enum = schema.Enum
			schema.Enum = nil
		}

	case TypeObject:
		for golangType.Kind() == reflect.Pointer {
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
			var structTypeName = oai.golangTypeToSchemaName(golangType)
			if oai.Components.Schemas.Get(structTypeName) == nil {
				if err = oai.addSchema(reflect.New(golangType).Interface()); err != nil {
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
					if err = oai.addSchema(golangTypeInstance); err != nil {
						return nil, err
					}
				}
				schemaRef.Ref = structTypeName
				schemaRef.Value = schema
				schemaRef.Description = schema.Description
			}
		}
	}
	oai.buildEnumValueXExtensions(typeId, schema, enumItems)
	return schemaRef, nil
}

func (oai *OpenApiV3) buildEnumValueXExtensions(typeId string, schema *Schema, enumItems []gtag.EnumItem) {
	if len(enumItems) == 0 || oai.Config.EnumXExtensionFunc == nil {
		return
	}
	var targetSchema = schema
	if schema.Type == TypeArray && schema.Items != nil && schema.Items.Value != nil {
		targetSchema = schema.Items.Value
	}
	if len(targetSchema.Enum) == 0 {
		return
	}
	if targetSchema.XExtensions == nil {
		targetSchema.XExtensions = make(XExtensions)
	}
	extensionMap := oai.Config.EnumXExtensionFunc(EnumXExtensionInput{
		TypeID: typeId,
		Items:  enumItems,
	})
	for extensionKey, extensionValue := range extensionMap {
		if extensionKey == "" {
			continue
		}
		if !gstr.HasPrefix(extensionKey, "x-") && !gstr.HasPrefix(extensionKey, "X-") {
			extensionKey = "x-" + extensionKey
		}
		targetSchema.XExtensions[extensionKey] = extensionValue
	}
}

func (r SchemaRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefAndDescToBytes(r.Ref, r.Description), nil
	}
	return json.Marshal(r.Value)
}

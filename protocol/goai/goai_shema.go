// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"reflect"
)

type Schemas map[string]SchemaRef

// Schema is specified by OpenAPI/Swagger 3.0 standard.
type Schema struct {
	OneOf                SchemaRefs     `json:"oneOf,omitempty"                yaml:"oneOf,omitempty"`
	AnyOf                SchemaRefs     `json:"anyOf,omitempty"                yaml:"anyOf,omitempty"`
	AllOf                SchemaRefs     `json:"allOf,omitempty"                yaml:"allOf,omitempty"`
	Not                  *SchemaRef     `json:"not,omitempty"                  yaml:"not,omitempty"`
	Type                 string         `json:"type,omitempty"                 yaml:"type,omitempty"`
	Title                string         `json:"title,omitempty"                yaml:"title,omitempty"`
	Format               string         `json:"format,omitempty"               yaml:"format,omitempty"`
	Description          string         `json:"description,omitempty"          yaml:"description,omitempty"`
	Enum                 []interface{}  `json:"enum,omitempty"                 yaml:"enum,omitempty"`
	Default              interface{}    `json:"default,omitempty"              yaml:"default,omitempty"`
	Example              interface{}    `json:"example,omitempty"              yaml:"example,omitempty"`
	ExternalDocs         *ExternalDocs  `json:"externalDocs,omitempty"         yaml:"externalDocs,omitempty"`
	UniqueItems          bool           `json:"uniqueItems,omitempty"          yaml:"uniqueItems,omitempty"`
	ExclusiveMin         bool           `json:"exclusiveMinimum,omitempty"     yaml:"exclusiveMinimum,omitempty"`
	ExclusiveMax         bool           `json:"exclusiveMaximum,omitempty"     yaml:"exclusiveMaximum,omitempty"`
	Nullable             bool           `json:"nullable,omitempty"             yaml:"nullable,omitempty"`
	ReadOnly             bool           `json:"readOnly,omitempty"             yaml:"readOnly,omitempty"`
	WriteOnly            bool           `json:"writeOnly,omitempty"            yaml:"writeOnly,omitempty"`
	AllowEmptyValue      bool           `json:"allowEmptyValue,omitempty"      yaml:"allowEmptyValue,omitempty"`
	XML                  interface{}    `json:"xml,omitempty"                  yaml:"xml,omitempty"`
	Deprecated           bool           `json:"deprecated,omitempty"           yaml:"deprecated,omitempty"`
	Min                  *float64       `json:"minimum,omitempty"              yaml:"minimum,omitempty"`
	Max                  *float64       `json:"maximum,omitempty"              yaml:"maximum,omitempty"`
	MultipleOf           *float64       `json:"multipleOf,omitempty"           yaml:"multipleOf,omitempty"`
	MinLength            uint64         `json:"minLength,omitempty"            yaml:"minLength,omitempty"`
	MaxLength            *uint64        `json:"maxLength,omitempty"            yaml:"maxLength,omitempty"`
	Pattern              string         `json:"pattern,omitempty"              yaml:"pattern,omitempty"`
	MinItems             uint64         `json:"minItems,omitempty"             yaml:"minItems,omitempty"`
	MaxItems             *uint64        `json:"maxItems,omitempty"             yaml:"maxItems,omitempty"`
	Items                *SchemaRef     `json:"items,omitempty"                yaml:"items,omitempty"`
	Required             []string       `json:"required,omitempty"             yaml:"required,omitempty"`
	Properties           Schemas        `json:"properties,omitempty"           yaml:"properties,omitempty"`
	MinProps             uint64         `json:"minProperties,omitempty"        yaml:"minProperties,omitempty"`
	MaxProps             *uint64        `json:"maxProperties,omitempty"        yaml:"maxProperties,omitempty"`
	AdditionalProperties *SchemaRef     `json:"additionalProperties,omitempty" yaml:"additionalProperties"`
	Discriminator        *Discriminator `json:"discriminator,omitempty"        yaml:"discriminator,omitempty"`
}

// Discriminator is specified by OpenAPI/Swagger standard version 3.0.
type Discriminator struct {
	PropertyName string            `json:"propertyName"      yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

func (oai *OpenApiV3) addSchema(object ...interface{}) error {
	for _, v := range object {
		if err := oai.doAddSchemaSingle(v); err != nil {
			return err
		}
	}
	return nil
}

func (oai *OpenApiV3) doAddSchemaSingle(object interface{}) error {
	if oai.Components.Schemas == nil {
		oai.Components.Schemas = map[string]SchemaRef{}
	}

	var (
		reflectType    = reflect.TypeOf(object)
		structTypeName = gstr.SubStrFromREx(reflectType.String(), ".")
	)

	// Already added.
	if _, ok := oai.Components.Schemas[structTypeName]; ok {
		return nil
	}
	// Take the holder first.
	oai.Components.Schemas[structTypeName] = SchemaRef{}

	structFields, _ := structs.Fields(structs.FieldsInput{
		Pointer:         object,
		RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
	})
	var (
		schema = &Schema{
			Properties: map[string]SchemaRef{},
		}
	)
	schema.Type = TypeObject
	for _, structField := range structFields {
		if !gstr.IsLetterUpper(structField.Name()[0]) {
			continue
		}
		var (
			fieldName = structField.Name()
		)
		if jsonName := structField.TagJsonName(); jsonName != "" {
			fieldName = jsonName
		}
		schemaRef, err := oai.newSchemaRefWithGolangType(
			structField.Type().Type,
			structField.TagMap(),
		)
		if err != nil {
			return err
		}
		schema.Properties[fieldName] = *schemaRef
	}
	oai.Components.Schemas[structTypeName] = SchemaRef{
		Ref:   "",
		Value: schema,
	}
	return nil
}

func (oai *OpenApiV3) golangTypeToOAIType(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return TypeString

	case reflect.Struct:
		return TypeObject

	case reflect.Slice, reflect.Array:

		return TypeArray

	case reflect.Bool:
		return TypeBoolean

	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return TypeNumber

	default:
		return TypeObject
	}
}

// golangTypeToOAIFormat converts and returns OpenAPI parameter format for given golang type `t`.
// Note that it does not return standard OpenAPI parameter format but custom format in golang type.
func (oai *OpenApiV3) golangTypeToOAIFormat(t reflect.Type) string {
	return t.String()
}

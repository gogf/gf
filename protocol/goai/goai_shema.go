// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gvalid"
)

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

// addSchema creates schemas with objects.
// Note that the `object` can be array alias like: `type Res []Item`.
func (oai *OpenApiV3) addSchema(object ...interface{}) error {
	for _, v := range object {
		if err := oai.doAddSchemaSingle(v); err != nil {
			return err
		}
	}
	return nil
}

func (oai *OpenApiV3) doAddSchemaSingle(object interface{}) error {
	if oai.Components.Schemas.refs == nil {
		oai.Components.Schemas.refs = gmap.NewListMap()
	}

	var (
		reflectType    = reflect.TypeOf(object)
		structTypeName = oai.golangTypeToSchemaName(reflectType)
	)

	// Already added.
	if oai.Components.Schemas.Get(structTypeName) != nil {
		return nil
	}
	// Take the holder first.
	oai.Components.Schemas.Set(structTypeName, SchemaRef{})

	schema, err := oai.structToSchema(object)
	if err != nil {
		return err
	}

	oai.Components.Schemas.Set(structTypeName, SchemaRef{
		Ref:   "",
		Value: schema,
	})
	return nil
}

// structToSchema converts and returns given struct object as Schema.
func (oai *OpenApiV3) structToSchema(object interface{}) (*Schema, error) {
	var (
		tagMap = gmeta.Data(object)
		schema = &Schema{
			Properties: createSchemas(),
		}
	)
	if len(tagMap) > 0 {
		if err := oai.tagMapToSchema(tagMap, schema); err != nil {
			return nil, err
		}
	}
	if schema.Type != "" && schema.Type != TypeObject {
		return schema, nil
	}
	// []struct.
	if utils.IsArray(object) {
		schema.Type = TypeArray
		subSchemaRef, err := oai.newSchemaRefWithGolangType(reflect.TypeOf(object).Elem(), nil)
		if err != nil {
			return nil, err
		}
		schema.Items = subSchemaRef
		return schema, nil
	}
	// struct.
	structFields, _ := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         object,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	schema.Type = TypeObject
	for _, structField := range structFields {
		if !gstr.IsLetterUpper(structField.Name()[0]) {
			continue
		}
		var fieldName = structField.Name()
		if jsonName := structField.TagJsonName(); jsonName != "" {
			fieldName = jsonName
		}
		schemaRef, err := oai.newSchemaRefWithGolangType(
			structField.Type().Type,
			structField.TagMap(),
		)
		if err != nil {
			return nil, err
		}
		schema.Properties.Set(fieldName, *schemaRef)
	}
	return schema, nil
}

func (oai *OpenApiV3) tagMapToSchema(tagMap map[string]string, schema *Schema) error {
	var mergedTagMap = oai.fileMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, schema); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Schema failed`)
	}
	// Validation info to OpenAPI schema pattern.
	for _, tag := range gvalid.GetTags() {
		if validationTagValue, ok := tagMap[tag]; ok {
			_, validationRule, _ := gvalid.ParseTagValue(validationTagValue)
			schema.Pattern = validationRule
			break
		}
	}
	return nil
}

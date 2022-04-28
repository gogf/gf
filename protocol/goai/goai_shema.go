// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"reflect"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gvalid"
)

// Schema is specified by OpenAPI/Swagger 3.0 standard.
type Schema struct {
	OneOf                SchemaRefs     `json:"oneOf,omitempty"`
	AnyOf                SchemaRefs     `json:"anyOf,omitempty"`
	AllOf                SchemaRefs     `json:"allOf,omitempty"`
	Not                  *SchemaRef     `json:"not,omitempty"`
	Type                 string         `json:"type,omitempty"`
	Title                string         `json:"title,omitempty"`
	Format               string         `json:"format,omitempty"`
	Description          string         `json:"description,omitempty"`
	Enum                 []interface{}  `json:"enum,omitempty"`
	Default              interface{}    `json:"default,omitempty"`
	Example              interface{}    `json:"example,omitempty"`
	ExternalDocs         *ExternalDocs  `json:"externalDocs,omitempty"`
	UniqueItems          bool           `json:"uniqueItems,omitempty"`
	ExclusiveMin         bool           `json:"exclusiveMinimum,omitempty"`
	ExclusiveMax         bool           `json:"exclusiveMaximum,omitempty"`
	Nullable             bool           `json:"nullable,omitempty"`
	ReadOnly             bool           `json:"readOnly,omitempty"`
	WriteOnly            bool           `json:"writeOnly,omitempty"`
	AllowEmptyValue      bool           `json:"allowEmptyValue,omitempty"`
	XML                  interface{}    `json:"xml,omitempty"`
	Deprecated           bool           `json:"deprecated,omitempty"`
	Min                  *float64       `json:"minimum,omitempty"`
	Max                  *float64       `json:"maximum,omitempty"`
	MultipleOf           *float64       `json:"multipleOf,omitempty"`
	MinLength            uint64         `json:"minLength,omitempty"`
	MaxLength            *uint64        `json:"maxLength,omitempty"`
	Pattern              string         `json:"pattern,omitempty"`
	MinItems             uint64         `json:"minItems,omitempty"`
	MaxItems             *uint64        `json:"maxItems,omitempty"`
	Items                *SchemaRef     `json:"items,omitempty"`
	Required             []string       `json:"required,omitempty"`
	Properties           Schemas        `json:"properties,omitempty"`
	MinProps             uint64         `json:"minProperties,omitempty"`
	MaxProps             *uint64        `json:"maxProperties,omitempty"`
	AdditionalProperties *SchemaRef     `json:"additionalProperties,omitempty"`
	Discriminator        *Discriminator `json:"discriminator,omitempty"`
	XExtensions          XExtensions    `json:"-"`
}

func (s Schema) MarshalJSON() ([]byte, error) {
	var (
		b   []byte
		m   map[string]json.RawMessage
		err error
	)
	type tempSchema Schema // To prevent JSON marshal recursion error.
	if b, err = json.Marshal(tempSchema(s)); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range s.XExtensions {
		if b, err = json.Marshal(v); err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

// Discriminator is specified by OpenAPI/Swagger standard version 3.0.
type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
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
			Properties:  createSchemas(),
			XExtensions: make(XExtensions),
		}
		ignoreProperties []interface{}
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

	schema.Properties.Iterator(func(key string, ref SchemaRef) bool {
		if ref.Value != nil && ref.Value.Pattern != "" {
			validationRuleSet := gset.NewStrSetFrom(gstr.Split(ref.Value.Pattern, "|"))
			if validationRuleSet.Contains(patternKeyForRequired) {
				schema.Required = append(schema.Required, key)
			}
		}
		if !isValidTag(key) {
			ignoreProperties = append(ignoreProperties, key)
		}
		return true
	})

	if len(ignoreProperties) > 0 {
		schema.Properties.Removes(ignoreProperties)
	}

	return schema, nil
}

func (oai *OpenApiV3) tagMapToSchema(tagMap map[string]string, schema *Schema) error {
	var mergedTagMap = oai.fileMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, schema); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Schema failed`)
	}
	oai.tagMapToXExtensions(mergedTagMap, schema.XExtensions)
	// Validation info to OpenAPI schema pattern.
	for _, tag := range gvalid.GetTags() {
		if validationTagValue, ok := tagMap[tag]; ok {
			_, validationRules, _ := gvalid.ParseTagValue(validationTagValue)
			schema.Pattern = validationRules
			// Enum checks.
			if len(schema.Enum) == 0 {
				for _, rule := range gstr.SplitAndTrim(validationRules, "|") {
					if gstr.HasPrefix(rule, patternKeyForIn) {
						schema.Enum = gconv.Interfaces(gstr.SplitAndTrim(rule[len(patternKeyForIn):], ","))
					}
				}
			}
			break
		}
	}
	return nil
}

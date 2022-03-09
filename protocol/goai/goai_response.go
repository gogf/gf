// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
)

// Response is specified by OpenAPI/Swagger 3.0 standard.
type Response struct {
	Description string  `json:"description"           yaml:"description"`
	Headers     Headers `json:"headers,omitempty"     yaml:"headers,omitempty"`
	Content     Content `json:"content,omitempty"     yaml:"content,omitempty"`
	Links       Links   `json:"links,omitempty"       yaml:"links,omitempty"`
}

// Responses is specified by OpenAPI/Swagger 3.0 standard.
type Responses map[string]ResponseRef

type ResponseRef struct {
	Ref   string
	Value *Response
}

func (r ResponseRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

type getResponseSchemaRefInput struct {
	BusinessStructName      string      // The business struct name.
	CommonResponseObject    interface{} // Common response object.
	CommonResponseDataField string      // Common response data field.
}

func (oai *OpenApiV3) getResponseSchemaRef(in getResponseSchemaRefInput) (*SchemaRef, error) {
	if in.CommonResponseObject == nil {
		return &SchemaRef{
			Ref: in.BusinessStructName,
		}, nil
	}

	var (
		dataFieldsPartsArray       = gstr.Split(in.CommonResponseDataField, ".")
		bizResponseStructSchemaRef = oai.Components.Schemas.Get(in.BusinessStructName)
		schema, err                = oai.structToSchema(in.CommonResponseObject)
	)
	if err != nil {
		return nil, err
	}
	if in.CommonResponseDataField == "" && bizResponseStructSchemaRef != nil {
		// Normal response.
		bizResponseStructSchemaRef.Value.Properties.Iterator(func(key string, ref SchemaRef) bool {
			schema.Properties.Set(key, ref)
			return true
		})
	} else {
		// Common response.
		structFields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         in.CommonResponseObject,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		for _, structField := range structFields {
			var fieldName = structField.Name()
			if jsonName := structField.TagJsonName(); jsonName != "" {
				fieldName = jsonName
			}
			switch len(dataFieldsPartsArray) {
			case 1:
				if structField.Name() == dataFieldsPartsArray[0] {
					if err = oai.tagMapToSchema(structField.TagMap(), bizResponseStructSchemaRef.Value); err != nil {
						return nil, err
					}
					schema.Properties.Set(fieldName, *bizResponseStructSchemaRef)
					break
				}
			default:
				// Recursively creating common response object schema.
				if structField.Name() == dataFieldsPartsArray[0] {
					var structFieldInstance = reflect.New(structField.Type().Type).Elem()
					schemaRef, err := oai.getResponseSchemaRef(getResponseSchemaRefInput{
						BusinessStructName:      in.BusinessStructName,
						CommonResponseObject:    structFieldInstance,
						CommonResponseDataField: gstr.Join(dataFieldsPartsArray[1:], "."),
					})
					if err != nil {
						return nil, err
					}
					schema.Properties.Set(fieldName, *schemaRef)
					break
				}
			}
		}
	}

	return &SchemaRef{
		Value: schema,
	}, nil
}

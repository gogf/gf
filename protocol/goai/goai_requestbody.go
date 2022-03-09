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

// RequestBody is specified by OpenAPI/Swagger 3.0 standard.
type RequestBody struct {
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool    `json:"required,omitempty"    yaml:"required,omitempty"`
	Content     Content `json:"content,omitempty"     yaml:"content,omitempty"`
}

type RequestBodyRef struct {
	Ref   string
	Value *RequestBody
}

func (r RequestBodyRef) MarshalJSON() ([]byte, error) {
	if r.Ref != "" {
		return formatRefToBytes(r.Ref), nil
	}
	return json.Marshal(r.Value)
}

type getRequestSchemaRefInput struct {
	BusinessStructName string
	RequestObject      interface{}
	RequestDataField   string
}

func (oai *OpenApiV3) getRequestSchemaRef(in getRequestSchemaRefInput) (*SchemaRef, error) {
	if oai.Config.CommonRequest == nil {
		return &SchemaRef{
			Ref: in.BusinessStructName,
		}, nil
	}

	var (
		dataFieldsPartsArray      = gstr.Split(in.RequestDataField, ".")
		bizRequestStructSchemaRef = oai.Components.Schemas.Get(in.BusinessStructName)
		schema, err               = oai.structToSchema(in.RequestObject)
	)
	if err != nil {
		return nil, err
	}
	if in.RequestDataField == "" && bizRequestStructSchemaRef != nil {
		// Normal request.
		bizRequestStructSchemaRef.Value.Properties.Iterator(func(key string, ref SchemaRef) bool {
			schema.Properties.Set(key, ref)
			return true
		})
	} else {
		// Common request.
		structFields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         in.RequestObject,
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
					if err = oai.tagMapToSchema(structField.TagMap(), bizRequestStructSchemaRef.Value); err != nil {
						return nil, err
					}
					schema.Properties.Set(fieldName, *bizRequestStructSchemaRef)
					break
				}
			default:
				if structField.Name() == dataFieldsPartsArray[0] {
					var structFieldInstance = reflect.New(structField.Type().Type).Elem()
					schemaRef, err := oai.getRequestSchemaRef(getRequestSchemaRefInput{
						BusinessStructName: in.BusinessStructName,
						RequestObject:      structFieldInstance,
						RequestDataField:   gstr.Join(dataFieldsPartsArray[1:], "."),
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

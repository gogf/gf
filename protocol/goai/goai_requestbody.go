// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"reflect"
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
		dataFieldsPartsArray                                      = gstr.Split(in.RequestDataField, ".")
		bizRequestStructSchemaRef, bizRequestStructSchemaRefExist = oai.Components.Schemas[in.BusinessStructName]
		schema, err                                               = oai.structToSchema(in.RequestObject)
	)
	if err != nil {
		return nil, err
	}
	if in.RequestDataField == "" && bizRequestStructSchemaRefExist {
		for k, v := range bizRequestStructSchemaRef.Value.Properties {
			schema.Properties[k] = v
		}
	} else {
		structFields, _ := structs.Fields(structs.FieldsInput{
			Pointer:         in.RequestObject,
			RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
		})
		for _, structField := range structFields {
			var (
				fieldName = structField.Name()
			)
			if jsonName := structField.TagJsonName(); jsonName != "" {
				fieldName = jsonName
			}
			switch len(dataFieldsPartsArray) {
			case 1:
				if structField.Name() == dataFieldsPartsArray[0] {
					schema.Properties[fieldName] = bizRequestStructSchemaRef
					break
				}
			default:
				if structField.Name() == dataFieldsPartsArray[0] {
					var (
						structFieldInstance = reflect.New(structField.Type().Type)
					)
					schemaRef, err := oai.getRequestSchemaRef(getRequestSchemaRefInput{
						BusinessStructName: in.BusinessStructName,
						RequestObject:      structFieldInstance,
						RequestDataField:   gstr.Join(dataFieldsPartsArray[1:], "."),
					})
					if err != nil {
						return nil, err
					}
					schema.Properties[fieldName] = *schemaRef
					break
				}
			}
		}
	}

	return &SchemaRef{
		Value: schema,
	}, nil
}

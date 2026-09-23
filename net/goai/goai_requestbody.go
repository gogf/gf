// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"reflect"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
)

// RequestBody is specified by OpenAPI/Swagger 3.0 standard.
type RequestBody struct {
	Description string  `json:"description,omitempty"`
	Required    bool    `json:"required,omitempty"`
	Content     Content `json:"content,omitempty"`
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
	RequestObject      any
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

	if bizRequestStructSchemaRef == nil {
		return &SchemaRef{
			Value: schema,
		}, nil
	}

	if in.RequestDataField == "" && bizRequestStructSchemaRef.Value != nil {
		// Append bizRequest.
		schema.Required = append(schema.Required, bizRequestStructSchemaRef.Value.Required...)

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

// getArrayRequestSchemaRef generates the OpenAPI schema reference for a JSON array request body,
// which is declared by the `type:"array"` tag in `g.Meta`. APIs using this tag accept a JSON
// array request body like `[{"id":1},{"id":2}]` instead of an object like `{"items":[{"id":1}]}`.
//
// The schema is generated from the first slice/array attribute of the request struct, for example:
//
//	type BatchChatReq struct {
//	    g.Meta   `mime:"application/json" method:"post" path:"/batch/chat" type:"array"`
//	    Messages []ChatMessage `json:"messages"`
//	}
//
// The generated OpenAPI definition is:
//
//	requestBody:
//	  content:
//	    application/json:
//	      schema:
//	        type: array
//	        items:
//	          $ref: '#/components/schemas/ChatMessage'
func (oai *OpenApiV3) getArrayRequestSchemaRef(requestObject any) (*SchemaRef, error) {
	structFields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         requestObject,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		return nil, err
	}
	// The request body definition is the first slice/array attribute of the request struct,
	// as only one attribute is able to receive the JSON array request body.
	for _, structField := range structFields {
		var golangType = structField.Type().Type
		if golangType.Kind() != reflect.Slice && golangType.Kind() != reflect.Array {
			continue
		}
		if structField.TagPriorityName() == "-" {
			continue
		}
		// It also recursively registers the schema of all the nested struct types in the
		// element type into Components.Schemas.
		elementSchemaRef, err := oai.newSchemaRefWithGolangType(golangType.Elem(), nil)
		if err != nil {
			return nil, err
		}
		return &SchemaRef{
			Value: &Schema{
				Type:  TypeArray,
				Items: elementSchemaRef,
			},
		}, nil
	}
	return nil, gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`there is no slice/array attribute in request struct "%s" for the type:"array" tag definition`,
		oai.golangTypeToSchemaName(reflect.TypeOf(requestObject)),
	)
}

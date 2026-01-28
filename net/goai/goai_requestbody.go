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

// getArrayRequestSchemaRef generates an OpenAPI schema reference for JSON array request bodies.
// This function supports the type:"array" tag in g.Meta for APIs that receive batch request formats.
//
// Implementation Steps:
//  1. Pre-register all nested struct types: Walk through the request struct and register any
//     struct fields and slice element types in Components.Schemas. This ensures proper schema
//     generation for nested structures.
//  2. Locate slice field: Use gstructs.Fields to find the first slice/array field in the struct.
//  3. Generate element schema: Create a schema reference for the slice's element type using
//     newSchemaRefWithGolangType, which handles nested structures recursively.
//  4. Return array schema: Construct a SchemaRef with Type="array" and Items pointing to the
//     element schema.
//
// Example OpenAPI Output:
//
//	requestBody:
//	  content:
//	    application/json:
//	      schema:
//	        type: array
//	        items:
//	          $ref: '#/components/schemas/ChatMessage'
//
// Related: ghttp_request_param_request.go::mergeBodyArrayToStruct
func (oai *OpenApiV3) getArrayRequestSchemaRef(requestObject any) (*SchemaRef, error) {
	// Step 1: Pre-register all nested struct types in Components.Schemas.
	// This is necessary because newSchemaRefWithGolangType only registers the direct
	// element type. We need to explicitly register nested structs to ensure proper
	// OpenAPI documentation generation.
	objectValue := reflect.ValueOf(requestObject)
	if objectValue.Kind() == reflect.Pointer {
		objectValue = objectValue.Elem()
	}
	if objectValue.Kind() == reflect.Struct {
		structType := objectValue.Type()
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldKind := field.Type.Kind()
			switch fieldKind {
			case reflect.Struct:
				// Register nested struct types
				_, _ = oai.newSchemaRefWithGolangType(field.Type, nil)
			case reflect.Slice, reflect.Array:
				// Register slice element types if they are structs
				elemType := field.Type.Elem()
				if elemType.Kind() == reflect.Struct {
					_, _ = oai.newSchemaRefWithGolangType(elemType, nil)
				}
			}
		}
	}

	// Step 2: Find the slice/array field in the struct.
	// We use gstructs.Fields to properly handle embedded structs.
	structFields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         requestObject,
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	if err != nil {
		return nil, err
	}

	var sliceField *gstructs.Field
	for _, field := range structFields {
		fieldValueType := field.Value.Type()
		if fieldValueType.Kind() == reflect.Slice || fieldValueType.Kind() == reflect.Array {
			sliceField = &field
			break
		}
	}

	// Step 3: Return empty array schema if no slice field found.
	if sliceField == nil {
		return &SchemaRef{
			Value: &Schema{
				Type:  "array",
				Items: &SchemaRef{},
			},
		}, nil
	}

	// Step 4: Get element type and generate schema reference.
	elementType := sliceField.Value.Type().Elem()
	elementSchemaRef, err := oai.newSchemaRefWithGolangType(elementType, nil)
	if err != nil {
		return nil, err
	}

	// Step 5: Return array schema with items referencing the element schema.
	return &SchemaRef{
		Value: &Schema{
			Type: "array",
			Items: &SchemaRef{
				Ref: elementSchemaRef.Ref,
			},
		},
	}, nil
}

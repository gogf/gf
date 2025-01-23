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
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

type ResponseRef struct {
	Ref   string
	Value *Response
}

// Responses is specified by OpenAPI/Swagger 3.0 standard.
type Responses map[string]ResponseRef

// object could be someObject.Interface()
// There may be some difference between someObject.Type() and reflect.TypeOf(object).
func (oai *OpenApiV3) getResponseFromObject(data interface{}, isDefault bool) (*Response, error) {
	var object interface{}
	enhancedResponse, isEnhanced := data.(EnhancedStatusType)
	if isEnhanced {
		object = enhancedResponse.Response
	} else {
		object = data
	}
	// Add object schema to oai
	if err := oai.addSchema(object); err != nil {
		return nil, err
	}
	var (
		metaMap  = gmeta.Data(object)
		response = &Response{
			Content:     map[string]MediaType{},
			XExtensions: make(XExtensions),
		}
	)
	if len(metaMap) > 0 {
		if err := oai.tagMapToResponse(metaMap, response); err != nil {
			return nil, err
		}
	}
	// Supported mime types of response.
	var (
		contentTypes = oai.Config.ReadContentTypes
		tagMimeValue = gmeta.Get(object, gtag.Mime).String()
		refInput     = getResponseSchemaRefInput{
			BusinessStructName:      oai.golangTypeToSchemaName(reflect.TypeOf(object)),
			CommonResponseObject:    oai.Config.CommonResponse,
			CommonResponseDataField: oai.Config.CommonResponseDataField,
		}
	)

	// If customized response mime type, it then ignores common response feature.
	if tagMimeValue != "" {
		contentTypes = gstr.SplitAndTrim(tagMimeValue, ",")
		refInput.CommonResponseObject = nil
		refInput.CommonResponseDataField = ""
	}

	// If it is not default status, check if it has any fields.
	// If so, it would override the common response.
	if !isDefault {
		fields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         object,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		if len(fields) > 0 {
			refInput.CommonResponseObject = nil
			refInput.CommonResponseDataField = ""
		}
	}

	// Generate response example from meta data.
	responseExamplePath := metaMap[gtag.ResponseExampleShort]
	if responseExamplePath == "" {
		responseExamplePath = metaMap[gtag.ResponseExample]
	}
	examples := make(Examples)
	if responseExamplePath != "" {
		if err := examples.applyExamplesFile(responseExamplePath); err != nil {
			return nil, err
		}
	}

	// Override examples from enhanced response.
	if isEnhanced {
		err := examples.applyExamplesData(enhancedResponse.Examples)
		if err != nil {
			return nil, err
		}
	}

	// Generate response schema from input.
	schemaRef, err := oai.getResponseSchemaRef(refInput)
	if err != nil {
		return nil, err
	}

	for _, contentType := range contentTypes {
		response.Content[contentType] = MediaType{
			Schema:   schemaRef,
			Examples: examples,
		}
	}
	return response, nil
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

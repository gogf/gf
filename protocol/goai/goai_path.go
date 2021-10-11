// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/structs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"reflect"
)

type Path struct {
	Ref         string     `json:"$ref,omitempty"        yaml:"$ref,omitempty"`
	Summary     string     `json:"summary,omitempty"     yaml:"summary,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Connect     *Operation `json:"connect,omitempty"     yaml:"connect,omitempty"`
	Delete      *Operation `json:"delete,omitempty"      yaml:"delete,omitempty"`
	Get         *Operation `json:"get,omitempty"         yaml:"get,omitempty"`
	Head        *Operation `json:"head,omitempty"        yaml:"head,omitempty"`
	Options     *Operation `json:"options,omitempty"     yaml:"options,omitempty"`
	Patch       *Operation `json:"patch,omitempty"       yaml:"patch,omitempty"`
	Post        *Operation `json:"post,omitempty"        yaml:"post,omitempty"`
	Put         *Operation `json:"put,omitempty"         yaml:"put,omitempty"`
	Trace       *Operation `json:"trace,omitempty"       yaml:"trace,omitempty"`
	Servers     Servers    `json:"servers,omitempty"     yaml:"servers,omitempty"`
	Parameters  Parameters `json:"parameters,omitempty"  yaml:"parameters,omitempty"`
}

// Paths is specified by OpenAPI/Swagger standard version 3.0.
type Paths map[string]Path

const (
	responseOkKey = `200`
)

type addPathInput struct {
	Path     string      // Precise route path.
	Prefix   string      // Route path prefix.
	Method   string      // Route method.
	Function interface{} // Uniformed function.
}

func (oai *OpenApiV3) addPath(in addPathInput) error {
	if oai.Paths == nil {
		oai.Paths = map[string]Path{}
	}

	var (
		reflectType = reflect.TypeOf(in.Function)
	)
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupported function "%s" for OpenAPI Path register, there should be input & output structures`,
			reflectType.String(),
		)
	}
	var (
		inputObject  reflect.Value
		outputObject reflect.Value
	)
	// Create instance according input/output types.
	if reflectType.In(1).Kind() == reflect.Ptr {
		inputObject = reflect.New(reflectType.In(1).Elem())
	} else {
		inputObject = reflect.New(reflectType.In(1).Elem()).Elem()
	}
	if reflectType.Out(0).Kind() == reflect.Ptr {
		outputObject = reflect.New(reflectType.Out(0).Elem()).Elem()
	} else {
		outputObject = reflect.New(reflectType.Out(0)).Elem()
	}
	for inputObject.Kind() == reflect.Ptr {
		inputObject = inputObject.Elem()
	}
	for outputObject.Kind() == reflect.Ptr {
		outputObject = outputObject.Elem()
	}
	var (
		path                 = Path{}
		inputMetaMap         = gmeta.Data(inputObject.Interface())
		outputMetaMap        = gmeta.Data(outputObject.Interface())
		isInputStructEmpty   = oai.doesStructHasNoFields(inputObject.Interface())
		inputStructTypeName  = golangTypeToSchemaName(inputObject.Type())
		outputStructTypeName = golangTypeToSchemaName(outputObject.Type())
		operation            = Operation{
			Responses: map[string]ResponseRef{},
		}
	)
	// Path check.
	if in.Path == "" {
		in.Path = gmeta.Get(inputObject.Interface(), TagNamePath).String()
		if in.Prefix != "" {
			in.Path = gstr.TrimRight(in.Prefix, "/") + "/" + gstr.TrimLeft(in.Path, "/")
		}
	}
	if in.Path == "" {
		return gerror.NewCodef(
			gcode.CodeMissingParameter,
			`missing necessary path parameter "%s" for input struct "%s", missing tag in attribute Meta?`,
			TagNamePath, inputStructTypeName,
		)
	}

	if v, ok := oai.Paths[in.Path]; ok {
		path = v
	}

	// Method check.
	if in.Method == "" {
		in.Method = gmeta.Get(inputObject.Interface(), TagNameMethod).String()
	}
	if in.Method == "" {
		return gerror.NewCodef(
			gcode.CodeMissingParameter,
			`missing necessary method parameter "%s" for input struct "%s", missing tag in attribute Meta?`,
			TagNameMethod, inputStructTypeName,
		)
	}

	if err := oai.addSchema(inputObject.Interface(), outputObject.Interface()); err != nil {
		return err
	}

	if len(inputMetaMap) > 0 {
		if err := gconv.Struct(inputMetaMap, &path); err != nil {
			return gerror.WrapCode(gcode.CodeInternalError, err, `mapping struct tags to Path failed`)
		}
		if err := gconv.Struct(inputMetaMap, &operation); err != nil {
			return gerror.WrapCode(gcode.CodeInternalError, err, `mapping struct tags to Operation failed`)
		}
	}

	// =================================================================================================================
	// Request.
	// =================================================================================================================
	if operation.RequestBody == nil {
		operation.RequestBody = &RequestBodyRef{}
	}
	if operation.RequestBody.Value == nil {
		var (
			requestBody = RequestBody{
				Required: true,
				Content:  map[string]MediaType{},
			}
		)
		// Supported mime types of request.
		var (
			contentTypes = oai.Config.ReadContentTypes
			tagMimeValue = gmeta.Get(inputObject.Interface(), TagNameMime).String()
		)
		if tagMimeValue != "" {
			contentTypes = gstr.SplitAndTrim(tagMimeValue, ",")
		}
		for _, v := range contentTypes {
			if isInputStructEmpty {
				requestBody.Content[v] = MediaType{}
			} else {
				schemaRef, err := oai.getRequestSchemaRef(getRequestSchemaRefInput{
					BusinessStructName: inputStructTypeName,
					RequestObject:      oai.Config.CommonRequest,
					RequestDataField:   oai.Config.CommonRequestDataField,
				})
				if err != nil {
					return err
				}
				requestBody.Content[v] = MediaType{
					Schema: schemaRef,
				}
			}
		}
		operation.RequestBody = &RequestBodyRef{
			Value: &requestBody,
		}
	}
	// Request parameters.
	structFields, _ := structs.Fields(structs.FieldsInput{
		Pointer:         inputObject.Interface(),
		RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
	})
	for _, structField := range structFields {
		if operation.Parameters == nil {
			operation.Parameters = []ParameterRef{}
		}
		parameterRef, err := oai.newParameterRefWithStructMethod(structField, in.Method)
		if err != nil {
			return err
		}
		if parameterRef != nil {
			operation.Parameters = append(operation.Parameters, *parameterRef)
		}
	}

	// =================================================================================================================
	// Response.
	// =================================================================================================================
	if _, ok := operation.Responses[responseOkKey]; !ok {
		var (
			response = Response{
				Content: map[string]MediaType{},
			}
		)
		if len(outputMetaMap) > 0 {
			if err := gconv.Struct(outputMetaMap, &response); err != nil {
				return gerror.WrapCode(gcode.CodeInternalError, err, `mapping struct tags to Response failed`)
			}
		}
		// Supported mime types of response.
		var (
			contentTypes = oai.Config.ReadContentTypes
			tagMimeValue = gmeta.Get(outputObject.Interface(), TagNameMime).String()
		)
		if tagMimeValue != "" {
			contentTypes = gstr.SplitAndTrim(tagMimeValue, ",")
		}
		for _, v := range contentTypes {
			schemaRef, err := oai.getResponseSchemaRef(getResponseSchemaRefInput{
				BusinessStructName: outputStructTypeName,
				ResponseObject:     oai.Config.CommonResponse,
				ResponseDataField:  oai.Config.CommonResponseDataField,
			})
			if err != nil {
				return err
			}
			response.Content[v] = MediaType{
				Schema: schemaRef,
			}
		}
		operation.Responses[responseOkKey] = ResponseRef{Value: &response}
	}

	// Assign to certain operation attribute.
	switch gstr.ToUpper(in.Method) {
	case HttpMethodGet:
		// GET operations cannot have a requestBody.
		operation.RequestBody = nil
		path.Get = &operation

	case HttpMethodPut:
		path.Put = &operation

	case HttpMethodPost:
		path.Post = &operation

	case HttpMethodDelete:
		// DELETE operations cannot have a requestBody.
		operation.RequestBody = nil
		path.Delete = &operation

	case HttpMethodConnect:
		// Nothing to do for Connect.

	case HttpMethodHead:
		path.Head = &operation

	case HttpMethodOptions:
		path.Options = &operation

	case HttpMethodPatch:
		path.Patch = &operation

	case HttpMethodTrace:
		path.Trace = &operation

	default:
		return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid method "%s"`, in.Method)
	}
	oai.Paths[in.Path] = path
	return nil
}

func (oai *OpenApiV3) doesStructHasNoFields(s interface{}) bool {
	structFields, _ := structs.Fields(structs.FieldsInput{
		Pointer:         s,
		RecursiveOption: structs.RecursiveOptionEmbeddedNoTag,
	})
	return len(structFields) == 0
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmeta"
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
	Path     string
	Method   string
	Function interface{}
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
		inputStructTypeName  = gstr.SubStrFromREx(inputObject.Type().String(), ".")
		outputStructTypeName = gstr.SubStrFromREx(outputObject.Type().String(), ".")
		operation            = Operation{
			Responses: map[string]ResponseRef{},
		}
	)
	if in.Path == "" {
		in.Path = gmeta.Get(inputObject.Interface(), TagNamePath).String()
	}
	if in.Path == "" {
		panic(gerror.NewCode(
			gcode.CodeMissingParameter,
			`missing necessary path parameter "%s" for struct "%s"`,
			TagNamePath, inputStructTypeName,
		))
	}
	if in.Method == "" {
		in.Method = gmeta.Get(inputObject.Interface(), TagNameMethod).String()
	}
	if in.Method == "" {
		panic(gerror.NewCode(
			gcode.CodeMissingParameter,
			`missing necessary method parameter "%s" for struct "%s"`,
			TagNamePath, inputStructTypeName,
		))
	}

	if err := oai.addSchema(inputObject.Interface(), outputObject.Interface()); err != nil {
		return err
	}

	if len(inputMetaMap) > 0 {
		if err := gconv.Struct(inputMetaMap, &path); err != nil {
			return gerror.WrapCodef(gcode.CodeInternalError, err, `mapping struct tags to Path failed`)
		}
		if err := gconv.Struct(inputMetaMap, &operation); err != nil {
			return gerror.WrapCodef(gcode.CodeInternalError, err, `mapping struct tags to Operation failed`)
		}
	}
	// Request.
	if operation.RequestBody.Value == nil {
		var (
			requestBody = RequestBody{
				Required: true,
				Content:  map[string]MediaType{},
			}
		)
		// Supported mime types of request.
		for _, v := range oai.Config.ReadContentTypes {
			requestBody.Content[v] = MediaType{
				Schema: &SchemaRef{
					Ref: inputStructTypeName,
				},
			}
		}
		operation.RequestBody = RequestBodyRef{
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
		parameterRef, err := oai.newParameterRefWithStructMethod(structField)
		if err != nil {
			return err
		}
		if parameterRef != nil {
			operation.Parameters = append(operation.Parameters, *parameterRef)
		}
	}

	// Response.
	if _, ok := operation.Responses[responseOkKey]; !ok {
		var (
			response = Response{
				Content: map[string]MediaType{},
			}
		)
		if len(outputMetaMap) > 0 {
			if err := gconv.Struct(outputMetaMap, &response); err != nil {
				return gerror.WrapCodef(gcode.CodeInternalError, err, `mapping struct tags to Response failed`)
			}
		}
		// Supported mime types of response.
		for _, v := range oai.Config.WriteContentTypes {
			response.Content[v] = MediaType{
				Schema: &SchemaRef{
					Ref: outputStructTypeName,
				},
			}
		}
		operation.Responses[responseOkKey] = ResponseRef{Value: &response}
	}

	// Assign to certain operation attribute.
	switch gstr.ToUpper(in.Method) {
	case HttpMethodGet:
		path.Get = &operation
	case HttpMethodPut:
		path.Put = &operation
	case HttpMethodPost:
		path.Post = &operation
	case HttpMethodDelete:
		path.Delete = &operation
	case HttpMethodConnect:
		path.Connect = &operation
	case HttpMethodHead:
		path.Head = &operation
	case HttpMethodOptions:
		path.Options = &operation
	case HttpMethodPatch:
		path.Patch = &operation
	case HttpMethodTrace:
		path.Trace = &operation
	default:
		panic(gerror.NewCode(gcode.CodeInvalidParameter, `invalid method "%s"`, in.Method))
	}
	oai.Paths[in.Path] = path
	return nil
}

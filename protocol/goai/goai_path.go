package goai

import (
	"context"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
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
	tagNamePath   = `path`
	tagNameMethod = `method`
	responseOkKey = `200`
)

type addPathInput struct {
	Path     string
	Method   string
	Function interface{}
}

func (oai *OpenApiV3) addPath(in addPathInput) {
	if oai.Paths == nil {
		oai.Paths = map[string]Path{}
	}

	var (
		ctx         = context.TODO()
		reflectType = reflect.TypeOf(in.Function)
	)
	if reflectType.NumIn() != 2 || reflectType.NumOut() != 2 {
		intlog.Printf(
			ctx,
			`unsupported function "%s" for OpenAPI Path register`,
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
		in.Path = gmeta.Get(inputObject.Interface(), tagNamePath).String()
	}
	if in.Path == "" {
		panic(gerror.NewCode(
			gcode.CodeInvalidParameter,
			`missing necessary path parameter "%s" for struct "%s"`,
			tagNamePath, inputStructTypeName,
		))
	}
	if in.Method == "" {
		in.Method = gmeta.Get(inputObject.Interface(), tagNameMethod).String()
	}
	if in.Method == "" {
		in.Method = oai.Config.DefaultMethod
	}
	oai.addSchema(inputObject.Interface(), outputObject.Interface())
	if len(inputMetaMap) > 0 {
		if err := gconv.Struct(inputMetaMap, &path); err != nil {
			intlog.Error(ctx, err)
		}
		if err := gconv.Struct(inputMetaMap, &operation); err != nil {
			intlog.Error(ctx, err)
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
	// Response.
	if _, ok := operation.Responses[responseOkKey]; !ok {
		var (
			response = Response{
				Content: map[string]MediaType{},
			}
		)
		if len(outputMetaMap) > 0 {
			if err := gconv.Struct(outputMetaMap, &response); err != nil {
				intlog.Error(ctx, err)
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
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package goai implements and provides document generating for OpenApi specification.
//
// https://editor.swagger.io/
package goai

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"reflect"
)

// OpenApiV3 is the structure defined from:
// https://swagger.io/specification/
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md
type OpenApiV3 struct {
	Config       Config                `json:"-"                      yaml:"-"`
	OpenAPI      string                `json:"openapi"                yaml:"openapi"`
	Components   Components            `json:"components,omitempty"   yaml:"components,omitempty"`
	Info         Info                  `json:"info"                   yaml:"info"`
	Paths        Paths                 `json:"paths"                  yaml:"paths"`
	Security     *SecurityRequirements `json:"security,omitempty"     yaml:"security,omitempty"`
	Servers      *Servers              `json:"servers,omitempty"      yaml:"servers,omitempty"`
	Tags         *Tags                 `json:"tags,omitempty"         yaml:"tags,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocs is specified by OpenAPI/Swagger standard version 3.0.
type ExternalDocs struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	HttpMethodAll     = `ALL`
	HttpMethodGet     = `GET`
	HttpMethodPut     = `PUT`
	HttpMethodPost    = `POST`
	HttpMethodDelete  = `DELETE`
	HttpMethodConnect = `CONNECT`
	HttpMethodHead    = `HEAD`
	HttpMethodOptions = `OPTIONS`
	HttpMethodPatch   = `PATCH`
	HttpMethodTrace   = `TRACE`
)

const (
	TypeNumber     = `number`
	TypeBoolean    = `boolean`
	TypeArray      = `array`
	TypeString     = `string`
	TypeObject     = `object`
	FormatInt32    = `int32`
	FormatInt64    = `int64`
	FormatDouble   = `double`
	FormatByte     = `byte`
	FormatBinary   = `binary`
	FormatDate     = `date`
	FormatDateTime = `date-time`
	FormatPassword = `password`
)

const (
	ParameterInHeader = `header`
	ParameterInPath   = `path`
	ParameterInQuery  = `query`
	ParameterInCookie = `cookie`
)

const (
	TagNamePath     = `path`
	TagNameMethod   = `method`
	TagNameMime     = `mime`
	TagNameValidate = `v`
)

var (
	defaultReadContentTypes  = []string{`application/json`}
	defaultWriteContentTypes = []string{`application/json`}
)

// New creates and returns a OpenApiV3 implements object.
func New() *OpenApiV3 {
	oai := &OpenApiV3{}
	oai.fillWithDefaultValue()
	return oai
}

// AddInput is the structured parameter for function OpenApiV3.Add.
type AddInput struct {
	Path   string      // Path specifies the custom path if this is not configured in Meta of struct tag.
	Prefix string      // Prefix specifies the custom route path prefix, which will be added with the path tag in Meta of struct tag.
	Method string      // Method specifies the custom HTTP method if this is not configured in Meta of struct tag.
	Object interface{} // Object can be an instance of struct or a route function.
}

// Add adds an instance of struct or a route function to OpenApiV3 definition implements.
func (oai *OpenApiV3) Add(in AddInput) error {
	var (
		reflectValue = reflect.ValueOf(in.Object)
	)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectValue.Kind() {
	case reflect.Struct:
		return oai.addSchema(in.Object)

	case reflect.Func:
		return oai.addPath(addPathInput{
			Path:     in.Path,
			Prefix:   in.Prefix,
			Method:   in.Method,
			Function: in.Object,
		})

	default:
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupported parameter type "%s", only struct/function type is supported`,
			reflect.TypeOf(in.Object).String(),
		)
	}
}

func (oai OpenApiV3) String() string {
	b, err := json.Marshal(oai)
	if err != nil {
		intlog.Error(context.TODO(), err)
	}
	return string(b)
}

func (oai *OpenApiV3) golangTypeToOAIType(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return TypeString

	case reflect.Struct:
		switch t.String() {
		case `time.Time`, `gtime.Time`:
			return TypeString
		}
		return TypeObject

	case reflect.Slice, reflect.Array:
		switch t.String() {
		case `[]uint8`:
			return TypeString
		}
		return TypeArray

	case reflect.Bool:
		return TypeBoolean

	case
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return TypeNumber

	default:
		return TypeObject
	}
}

// golangTypeToOAIFormat converts and returns OpenAPI parameter format for given golang type `t`.
// Note that it does not return standard OpenAPI parameter format but custom format in golang type.
func (oai *OpenApiV3) golangTypeToOAIFormat(t reflect.Type) string {
	format := t.String()
	switch gstr.TrimLeft(format, "*") {
	case `[]uint8`:
		return FormatBinary

	default:
		return format
	}
}

func formatRefToBytes(ref string) []byte {
	return []byte(fmt.Sprintf(`{"$ref":"#/components/schemas/%s"}`, ref))
}

func golangTypeToSchemaName(t reflect.Type) string {
	var (
		s = gstr.TrimLeft(t.String(), "*")
	)
	if pkgPath := t.PkgPath(); pkgPath != "" && pkgPath != "." {
		s = gstr.Replace(t.PkgPath(), `/`, `.`) + gstr.SubStrFrom(s, ".")
	}
	s = gstr.ReplaceByMap(s, map[string]string{
		` `: ``,
		`{`: ``,
		`}`: ``,
	})
	return s
}

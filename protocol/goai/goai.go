package goai

import (
	"context"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/internal/json"
	"reflect"
)

// OpenApiV3 is the structure defined from:
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

type Config struct {
	DefaultMethod     string
	ReadContentTypes  []string
	WriteContentTypes []string
}

// ExternalDocs is specified by OpenAPI/Swagger standard version 3.0.
type ExternalDocs struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
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
	TypeInteger    = `integer`
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
	defaultMethod  = `POST`
)

const (
	defaultReadContentType  = `application/json`
	defaultWriteContentType = `application/json`
)

func New() *OpenApiV3 {
	oai := &OpenApiV3{}
	oai.fillWithDefaultValue()
	return oai
}

type AddInput struct {
	Path   string
	Method string
	Object interface{}
}

func (oai *OpenApiV3) Add(in AddInput) {
	var (
		reflectValue = reflect.ValueOf(in.Object)
	)
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	switch reflectValue.Kind() {
	case reflect.Struct:
		oai.addSchema(in.Object)

	case reflect.Func:
		oai.addPath(addPathInput{
			Path:     in.Path,
			Method:   in.Method,
			Function: in.Object,
		})

	default:
		panic(gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupported parameter type "%s", only struct/function type is supported`,
			reflect.TypeOf(in.Object).String(),
		))
	}
}

func (oai *OpenApiV3) fillWithDefaultValue() {
	if oai.OpenAPI == "" {
		oai.OpenAPI = `3.0.0`
	}
	if oai.Config.DefaultMethod == "" {
		oai.Config.DefaultMethod = defaultMethod
	}
	if len(oai.Config.ReadContentTypes) == 0 {
		oai.Config.ReadContentTypes = []string{defaultReadContentType}
	}
	if len(oai.Config.WriteContentTypes) == 0 {
		oai.Config.WriteContentTypes = []string{defaultWriteContentType}
	}
}

func (oai OpenApiV3) String() string {
	b, err := json.Marshal(oai)
	if err != nil {
		intlog.Error(context.TODO(), err)
	}
	return string(b)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"net/http"
	"reflect"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
)

// Path is specified by OpenAPI/Swagger standard version 3.0.
type Path struct {
	Ref         string      `json:"$ref,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Connect     *Operation  `json:"connect,omitempty"`
	Delete      *Operation  `json:"delete,omitempty"`
	Get         *Operation  `json:"get,omitempty"`
	Head        *Operation  `json:"head,omitempty"`
	Options     *Operation  `json:"options,omitempty"`
	Patch       *Operation  `json:"patch,omitempty"`
	Post        *Operation  `json:"post,omitempty"`
	Put         *Operation  `json:"put,omitempty"`
	Trace       *Operation  `json:"trace,omitempty"`
	Servers     Servers     `json:"servers,omitempty"`
	Parameters  Parameters  `json:"parameters,omitempty"`
	XExtensions XExtensions `json:"-"`
}

// Paths are specified by OpenAPI/Swagger standard version 3.0.
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

	var reflectType = reflect.TypeOf(in.Function)
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
		inputObject = reflect.New(reflectType.In(1).Elem()).Elem()
	} else {
		inputObject = reflect.New(reflectType.In(1)).Elem()
	}
	if reflectType.Out(0).Kind() == reflect.Ptr {
		outputObject = reflect.New(reflectType.Out(0).Elem()).Elem()
	} else {
		outputObject = reflect.New(reflectType.Out(0)).Elem()
	}

	var (
		mime                string
		path                = Path{XExtensions: make(XExtensions)}
		inputMetaMap        = gmeta.Data(inputObject.Interface())
		outputMetaMap       = gmeta.Data(outputObject.Interface())
		isInputStructEmpty  = oai.doesStructHasNoFields(inputObject.Interface())
		inputStructTypeName = oai.golangTypeToSchemaName(inputObject.Type())
		operation           = Operation{
			Responses:   map[string]ResponseRef{},
			XExtensions: make(XExtensions),
		}
		seRequirement = SecurityRequirement{}
	)
	// Path check.
	if in.Path == "" {
		in.Path = gmeta.Get(inputObject.Interface(), gtag.Path).String()
		if in.Prefix != "" {
			in.Path = gstr.TrimRight(in.Prefix, "/") + "/" + gstr.TrimLeft(in.Path, "/")
		}
	}
	if in.Path == "" {
		return gerror.NewCodef(
			gcode.CodeMissingParameter,
			`missing necessary path parameter "%s" for input struct "%s", missing tag in attribute Meta?`,
			gtag.Path, inputStructTypeName,
		)
	}

	if v, ok := oai.Paths[in.Path]; ok {
		path = v
	}

	// Method check.
	if in.Method == "" {
		in.Method = gmeta.Get(inputObject.Interface(), gtag.Method).String()
	}
	if in.Method == "" {
		return gerror.NewCodef(
			gcode.CodeMissingParameter,
			`missing necessary method parameter "%s" for input struct "%s", missing tag in attribute Meta?`,
			gtag.Method, inputStructTypeName,
		)
	}

	if err := oai.addSchema(inputObject.Interface()); err != nil {
		return err
	}

	if len(inputMetaMap) > 0 {
		// Path and Operation are not the same thing, so it is necessary to copy a Meta for Path from Operation and edit it.
		// And you know, we set the Summary and Description for Operation, not for Path, so we need to remove them.
		inputMetaMapForPath := gmap.NewStrStrMapFrom(inputMetaMap).Clone()
		inputMetaMapForPath.Removes([]string{
			gtag.SummaryShort,
			gtag.SummaryShort2,
			gtag.Summary,
			gtag.DescriptionShort,
			gtag.DescriptionShort2,
			gtag.Description,
		})
		if err := oai.tagMapToPath(inputMetaMapForPath.Map(), &path); err != nil {
			return err
		}

		if err := oai.tagMapToOperation(inputMetaMap, &operation); err != nil {
			return err
		}
		// Allowed request mime.
		if mime = inputMetaMap[gtag.Mime]; mime == "" {
			mime = inputMetaMap[gtag.Consumes]
		}
	}

	// path security
	// note: the security schema type only support http and apiKey;not support oauth2 and openIdConnect.
	// multi schema separate with comma, e.g. `security: apiKey1,apiKey2`
	TagNameSecurity := gmeta.Get(inputObject.Interface(), gtag.Security).String()
	securities := gstr.SplitAndTrim(TagNameSecurity, ",")
	for _, sec := range securities {
		seRequirement[sec] = []string{}
	}
	if len(securities) > 0 {
		operation.Security = &SecurityRequirements{seRequirement}
	}

	// =================================================================================================================
	// Request Parameter.
	// =================================================================================================================
	structFields, _ := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         inputObject.Interface(),
		RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
	})
	for _, structField := range structFields {
		if operation.Parameters == nil {
			operation.Parameters = []ParameterRef{}
		}
		parameterRef, err := oai.newParameterRefWithStructMethod(structField, in.Path, in.Method)
		if err != nil {
			return err
		}
		if parameterRef != nil {
			operation.Parameters = append(operation.Parameters, *parameterRef)
		}
	}

	// =================================================================================================================
	// Request Body.
	// =================================================================================================================
	if operation.RequestBody == nil {
		operation.RequestBody = &RequestBodyRef{}
	}
	if operation.RequestBody.Value == nil {
		var (
			requestBody = RequestBody{
				Content: map[string]MediaType{},
			}
		)
		// Supported mime types of request.
		var (
			contentTypes     = oai.Config.ReadContentTypes
			tagMimeValue     = gmeta.Get(inputObject.Interface(), gtag.Mime).String()
			tagRequiredValue = gmeta.Get(inputObject.Interface(), gtag.Required).Bool()
		)
		requestBody.Required = tagRequiredValue
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

	// =================================================================================================================
	// Default Response.
	// =================================================================================================================
	status := responseOkKey
	if statusValue, ok := outputMetaMap[gtag.Status]; ok {
		statusCode := gconv.Int(statusValue)
		if statusCode < 100 || statusCode >= 600 {
			return gerror.Newf("Invalid HTTP status code: %s", statusValue)
		}
		status = statusValue
	}
	if _, ok := operation.Responses[status]; !ok {
		response, err := oai.getResponseFromObject(outputObject.Interface(), true)
		if err != nil {
			return err
		}
		operation.Responses[status] = ResponseRef{Value: response}
	}

	// =================================================================================================================
	// Other Responses.
	// =================================================================================================================
	if enhancedResponse, ok := outputObject.Interface().(ResponseStatusDef); ok {
		for statusCode, data := range enhancedResponse.ResponseStatusMap() {
			if statusCode < 100 || statusCode >= 600 {
				return gerror.Newf("Invalid HTTP status code: %d", statusCode)
			}
			if data == nil {
				continue
			}
			status := gconv.String(statusCode)
			if _, ok := operation.Responses[status]; !ok {
				response, err := oai.getResponseFromObject(data, false)
				if err != nil {
					return err
				}
				operation.Responses[status] = ResponseRef{Value: response}
			}
		}
	}

	// Remove operation body duplicated properties.
	oai.removeOperationDuplicatedProperties(operation)

	// Assign to certain operation attribute.
	switch gstr.ToUpper(in.Method) {
	case http.MethodGet:
		// GET operations cannot have a requestBody.
		operation.RequestBody = nil
		path.Get = &operation

	case http.MethodPut:
		path.Put = &operation

	case http.MethodPost:
		path.Post = &operation

	case http.MethodDelete:
		// DELETE operations cannot have a requestBody.
		operation.RequestBody = nil
		path.Delete = &operation

	case http.MethodConnect:
		// Nothing to do for Connect.

	case http.MethodHead:
		path.Head = &operation

	case http.MethodOptions:
		path.Options = &operation

	case http.MethodPatch:
		path.Patch = &operation

	case http.MethodTrace:
		path.Trace = &operation

	default:
		return gerror.NewCodef(gcode.CodeInvalidParameter, `invalid method "%s"`, in.Method)
	}
	oai.Paths[in.Path] = path
	return nil
}

func (oai *OpenApiV3) removeOperationDuplicatedProperties(operation Operation) {
	if len(operation.Parameters) == 0 {
		// Nothing to do.
		return
	}

	var (
		duplicatedParameterNames []interface{}
		dataField                string
	)

	for _, parameter := range operation.Parameters {
		duplicatedParameterNames = append(duplicatedParameterNames, parameter.Value.Name)
	}

	// Check operation request body have common request data field.
	dataFields := gstr.Split(oai.Config.CommonRequestDataField, ".")
	if len(dataFields) > 0 && dataFields[0] != "" {
		dataField = dataFields[0]
	}

	for _, requestBodyContent := range operation.RequestBody.Value.Content {
		// Check request body schema
		if requestBodyContent.Schema == nil {
			continue
		}

		// Check request body schema ref.
		if requestBodyContent.Schema.Ref != "" {
			if schema := oai.Components.Schemas.Get(requestBodyContent.Schema.Ref); schema != nil {
				newSchema := schema.Value.Clone()
				requestBodyContent.Schema.Ref = ""
				requestBodyContent.Schema.Value = newSchema
				newSchema.Required = oai.removeItemsFromArray(newSchema.Required, duplicatedParameterNames)
				newSchema.Properties.Removes(duplicatedParameterNames)
				continue
			}
		}

		// Check the Value public field for the request body.
		if commonRequest := requestBodyContent.Schema.Value.Properties.Get(dataField); commonRequest != nil {
			commonRequest.Value.Required = oai.removeItemsFromArray(commonRequest.Value.Required, duplicatedParameterNames)
			commonRequest.Value.Properties.Removes(duplicatedParameterNames)
			continue
		}

		// Check request body schema value.
		if requestBodyContent.Schema.Value != nil {
			requestBodyContent.Schema.Value.Required = oai.removeItemsFromArray(requestBodyContent.Schema.Value.Required, duplicatedParameterNames)
			requestBodyContent.Schema.Value.Properties.Removes(duplicatedParameterNames)
			continue
		}
	}
}

func (oai *OpenApiV3) removeItemsFromArray(array []string, items []interface{}) []string {
	arr := garray.NewStrArrayFrom(array)
	for _, item := range items {
		if value, ok := item.(string); ok {
			arr.RemoveValue(value)
		}
	}
	return arr.Slice()
}

func (oai *OpenApiV3) doesStructHasNoFields(s interface{}) bool {
	return reflect.TypeOf(s).NumField() == 0
}

func (oai *OpenApiV3) tagMapToPath(tagMap map[string]string, path *Path) error {
	var mergedTagMap = oai.fillMapWithShortTags(tagMap)
	if err := gconv.Struct(mergedTagMap, path); err != nil {
		return gerror.Wrap(err, `mapping struct tags to Path failed`)
	}
	oai.tagMapToXExtensions(mergedTagMap, path.XExtensions)
	return nil
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (p Path) MarshalJSON() ([]byte, error) {
	var (
		b   []byte
		m   map[string]json.RawMessage
		err error
	)
	type tempPath Path // To prevent JSON marshal recursion error.
	if b, err = json.Marshal(tempPath(p)); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	for k, v := range p.XExtensions {
		if b, err = json.Marshal(v); err != nil {
			return nil, err
		}
		m[k] = b
	}
	return json.Marshal(m)
}

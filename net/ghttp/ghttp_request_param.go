// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/encoding/gxml"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gvalid"
	"io/ioutil"
	"mime/multipart"
	"reflect"
	"strings"
)

const (
	parseTypeRequest = 0
	parseTypeQuery   = 1
	parseTypeForm    = 2
)

var (
	// xmlHeaderBytes is the most common XML format header.
	xmlHeaderBytes = []byte("<?xml")
)

// Parse is the most commonly used function, which converts request parameters to struct or struct
// slice. It also automatically validates the struct or every element of the struct slice according
// to the validation tag of the struct.
//
// The parameter <pointer> can be type of: *struct/**struct/*[]struct/*[]*struct.
//
// It supports single and multiple struct convertion:
// 1. Single struct, post content like: {"id":1, "name":"john"} or ?id=1&name=john
// 2. Multiple struct, post content like: [{"id":1, "name":"john"}, {"id":, "name":"smith"}]
//
// TODO: Improve the performance by reducing duplicated reflect usage on the same variable across packages.
func (r *Request) Parse(pointer interface{}) error {
	return r.doParse(pointer, parseTypeRequest)
}

// ParseQuery performs like function Parse, but only parses the query parameters.
func (r *Request) ParseQuery(pointer interface{}) error {
	return r.doParse(pointer, parseTypeQuery)
}

// ParseForm performs like function Parse, but only parses the form parameters or the body content.
func (r *Request) ParseForm(pointer interface{}) error {
	return r.doParse(pointer, parseTypeForm)
}

// doParse parses the request data to struct/structs according to request type.
func (r *Request) doParse(pointer interface{}, requestType int) error {
	var (
		reflectVal1  = reflect.ValueOf(pointer)
		reflectKind1 = reflectVal1.Kind()
	)
	if reflectKind1 != reflect.Ptr {
		return fmt.Errorf(
			"parameter should be type of *struct/**struct/*[]struct/*[]*struct, but got: %v",
			reflectKind1,
		)
	}
	var (
		reflectVal2  = reflectVal1.Elem()
		reflectKind2 = reflectVal2.Kind()
	)
	switch reflectKind2 {
	// Single struct, post content like:
	// 1. {"id":1, "name":"john"}
	// 2. ?id=1&name=john
	case reflect.Ptr, reflect.Struct:
		// Converting.
		switch requestType {
		case parseTypeQuery:
			if err := r.GetQueryStruct(pointer); err != nil {
				return err
			}
		case parseTypeForm:
			if err := r.GetFormStruct(pointer); err != nil {
				return err
			}
		default:
			if err := r.GetStruct(pointer); err != nil {
				return err
			}
		}

		// Validation.
		if err := gvalid.CheckStruct(pointer, nil); err != nil {
			return err
		}

	// Multiple struct, it only supports JSON type post content like:
	// [{"id":1, "name":"john"}, {"id":, "name":"smith"}]
	case reflect.Array, reflect.Slice:
		// If struct slice conversion, it might post JSON/XML content,
		// so it uses gjson for the conversion.
		j, err := gjson.LoadContent(r.GetBody())
		if err != nil {
			return err
		}
		if err := j.GetStructs(".", pointer); err != nil {
			return err
		}
		for i := 0; i < reflectVal2.Len(); i++ {
			if err := gvalid.CheckStruct(reflectVal2.Index(i), nil); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get is alias of GetRequest, which is one of the most commonly used functions for
// retrieving parameter.
// See r.GetRequest.
func (r *Request) Get(key string, def ...interface{}) interface{} {
	return r.GetRequest(key, def...)
}

// GetVar is alis of GetRequestVar.
// See GetRequestVar.
func (r *Request) GetVar(key string, def ...interface{}) *gvar.Var {
	return r.GetRequestVar(key, def...)
}

// GetRaw is alias of GetBody.
// See GetBody.
// Deprecated.
func (r *Request) GetRaw() []byte {
	return r.GetBody()
}

// GetRawString is alias of GetBodyString.
// See GetBodyString.
// Deprecated.
func (r *Request) GetRawString() string {
	return r.GetBodyString()
}

// GetBody retrieves and returns request body content as bytes.
// It can be called multiple times retrieving the same body content.
func (r *Request) GetBody() []byte {
	if r.bodyContent == nil {
		r.bodyContent, _ = ioutil.ReadAll(r.Body)
		r.Body = utils.NewReadCloser(r.bodyContent, true)
	}
	return r.bodyContent
}

// GetBodyString retrieves and returns request body content as string.
// It can be called multiple times retrieving the same body content.
func (r *Request) GetBodyString() string {
	return gconv.UnsafeBytesToStr(r.GetBody())
}

// GetJson parses current request content as JSON format, and returns the JSON object.
// Note that the request content is read from request BODY, not from any field of FORM.
func (r *Request) GetJson() (*gjson.Json, error) {
	return gjson.LoadJson(r.GetBody())
}

// GetString is an alias and convenient function for GetRequestString.
// See GetRequestString.
func (r *Request) GetString(key string, def ...interface{}) string {
	return r.GetRequestString(key, def...)
}

// GetBool is an alias and convenient function for GetRequestBool.
// See GetRequestBool.
func (r *Request) GetBool(key string, def ...interface{}) bool {
	return r.GetRequestBool(key, def...)
}

// GetInt is an alias and convenient function for GetRequestInt.
// See GetRequestInt.
func (r *Request) GetInt(key string, def ...interface{}) int {
	return r.GetRequestInt(key, def...)
}

// GetInt32 is an alias and convenient function for GetRequestInt32.
// See GetRequestInt32.
func (r *Request) GetInt32(key string, def ...interface{}) int32 {
	return r.GetRequestInt32(key, def...)
}

// GetInt64 is an alias and convenient function for GetRequestInt64.
// See GetRequestInt64.
func (r *Request) GetInt64(key string, def ...interface{}) int64 {
	return r.GetRequestInt64(key, def...)
}

// GetInts is an alias and convenient function for GetRequestInts.
// See GetRequestInts.
func (r *Request) GetInts(key string, def ...interface{}) []int {
	return r.GetRequestInts(key, def...)
}

// GetUint is an alias and convenient function for GetRequestUint.
// See GetRequestUint.
func (r *Request) GetUint(key string, def ...interface{}) uint {
	return r.GetRequestUint(key, def...)
}

// GetUint32 is an alias and convenient function for GetRequestUint32.
// See GetRequestUint32.
func (r *Request) GetUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestUint32(key, def...)
}

// GetUint64 is an alias and convenient function for GetRequestUint64.
// See GetRequestUint64.
func (r *Request) GetUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestUint64(key, def...)
}

// GetFloat32 is an alias and convenient function for GetRequestFloat32.
// See GetRequestFloat32.
func (r *Request) GetFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestFloat32(key, def...)
}

// GetFloat64 is an alias and convenient function for GetRequestFloat64.
// See GetRequestFloat64.
func (r *Request) GetFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestFloat64(key, def...)
}

// GetFloats is an alias and convenient function for GetRequestFloats.
// See GetRequestFloats.
func (r *Request) GetFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestFloats(key, def...)
}

// GetArray is an alias and convenient function for GetRequestArray.
// See GetRequestArray.
func (r *Request) GetArray(key string, def ...interface{}) []string {
	return r.GetRequestArray(key, def...)
}

// GetStrings is an alias and convenient function for GetRequestStrings.
// See GetRequestStrings.
func (r *Request) GetStrings(key string, def ...interface{}) []string {
	return r.GetRequestStrings(key, def...)
}

// GetInterfaces is an alias and convenient function for GetRequestInterfaces.
// See GetRequestInterfaces.
func (r *Request) GetInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestInterfaces(key, def...)
}

// GetMap is an alias and convenient function for GetRequestMap.
// See GetRequestMap.
func (r *Request) GetMap(def ...map[string]interface{}) map[string]interface{} {
	return r.GetRequestMap(def...)
}

// GetMapStrStr is an alias and convenient function for GetRequestMapStrStr.
// See GetRequestMapStrStr.
func (r *Request) GetMapStrStr(def ...map[string]interface{}) map[string]string {
	return r.GetRequestMapStrStr(def...)
}

// GetStruct is an alias and convenient function for GetRequestStruct.
// See GetRequestStruct.
func (r *Request) GetStruct(pointer interface{}, mapping ...map[string]string) error {
	return r.GetRequestStruct(pointer, mapping...)
}

// parseQuery parses query string into r.queryMap.
func (r *Request) parseQuery() {
	if r.parsedQuery {
		return
	}
	r.parsedQuery = true
	if r.URL.RawQuery != "" {
		var err error
		r.queryMap, err = gstr.Parse(r.URL.RawQuery)
		if err != nil {
			panic(err)
		}
	}
}

// parseBody parses the request raw data into r.rawMap.
// Note that it also supports JSON data from client request.
func (r *Request) parseBody() {
	if r.parsedBody {
		return
	}
	r.parsedBody = true
	// There's no data posted.
	if r.ContentLength == 0 {
		return
	}
	if body := r.GetBody(); len(body) > 0 {
		// Trim space/new line characters.
		body = bytes.TrimSpace(body)
		// JSON format checks.
		if body[0] == '{' && body[len(body)-1] == '}' {
			_ = json.Unmarshal(body, &r.bodyMap)
		}
		// XML format checks.
		if len(body) > 5 && bytes.EqualFold(body[:5], xmlHeaderBytes) {
			r.bodyMap, _ = gxml.DecodeWithoutRoot(body)
		}
		if body[0] == '<' && body[len(body)-1] == '>' {
			r.bodyMap, _ = gxml.DecodeWithoutRoot(body)
		}
		// Default parameters decoding.
		if r.bodyMap == nil {
			r.bodyMap, _ = gstr.Parse(r.GetBodyString())
		}
	}
}

// parseForm parses the request form for HTTP method PUT, POST, PATCH.
// The form data is pared into r.formMap.
//
// Note that if the form was parsed firstly, the request body would be cleared and empty.
func (r *Request) parseForm() {
	if r.parsedForm {
		return
	}
	r.parsedForm = true
	// There's no data posted.
	if r.ContentLength == 0 {
		return
	}
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		var err error
		if gstr.Contains(contentType, "multipart/") {
			// multipart/form-data, multipart/mixed
			if err = r.ParseMultipartForm(r.Server.config.FormParsingMemory); err != nil {
				panic(err)
			}
		} else if gstr.Contains(contentType, "form") {
			// application/x-www-form-urlencoded
			if err = r.Request.ParseForm(); err != nil {
				panic(err)
			}
		}
		if len(r.PostForm) > 0 {
			// Re-parse the form data using united parsing way.
			params := ""
			for name, values := range r.PostForm {
				// Invalid parameter name.
				// Only allow chars of: '\w', '[', ']', '-'.
				if !gregex.IsMatchString(`^[\w\-\[\]]+$`, name) && len(r.PostForm) == 1 {
					// It might be JSON/XML content.
					if s := gstr.Trim(name + strings.Join(values, " ")); len(s) > 0 {
						if s[0] == '{' && s[len(s)-1] == '}' || s[0] == '<' && s[len(s)-1] == '>' {
							r.bodyContent = gconv.UnsafeStrToBytes(s)
							params = ""
							break
						}
					}
				}
				if len(values) == 1 {
					if len(params) > 0 {
						params += "&"
					}
					params += name + "=" + gurl.Encode(values[0])
				} else {
					if len(name) > 2 && name[len(name)-2:] == "[]" {
						name = name[:len(name)-2]
						for _, v := range values {
							if len(params) > 0 {
								params += "&"
							}
							params += name + "[]=" + gurl.Encode(v)
						}
					} else {
						if len(params) > 0 {
							params += "&"
						}
						params += name + "=" + gurl.Encode(values[len(values)-1])
					}
				}
			}
			if params != "" {
				if r.formMap, err = gstr.Parse(params); err != nil {
					panic(err)
				}
			}
		}
	}
	// It parses the request body without checking the Content-Type.
	if r.formMap == nil {
		if r.Method != "GET" {
			r.parseBody()
		}
		if len(r.bodyMap) > 0 {
			r.formMap = r.bodyMap
		}
	}
}

// GetMultipartForm parses and returns the form as multipart form.
func (r *Request) GetMultipartForm() *multipart.Form {
	r.parseForm()
	return r.MultipartForm
}

// GetMultipartFiles parses and returns the post files array.
// Note that the request form should be type of multipart.
func (r *Request) GetMultipartFiles(name string) []*multipart.FileHeader {
	form := r.GetMultipartForm()
	if form == nil {
		return nil
	}
	if v := form.File[name]; len(v) > 0 {
		return v
	}
	// Support "name[]" as array parameter.
	if v := form.File[name+"[]"]; len(v) > 0 {
		return v
	}
	// Support "name[0]","name[1]","name[2]", etc. as array parameter.
	key := ""
	files := make([]*multipart.FileHeader, 0)
	for i := 0; ; i++ {
		key = fmt.Sprintf(`%s[%d]`, name, i)
		if v := form.File[key]; len(v) > 0 {
			files = append(files, v[0])
		} else {
			break
		}
	}
	if len(files) > 0 {
		return files
	}
	return nil
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// GetRequest retrieves and returns the parameter named `key` passed from the client and
// custom params as any, no matter what HTTP method the client is using. The
// parameter `def` specifies the default value if the `key` does not exist.
//
// GetRequest is one of the most commonly used functions for retrieving parameters.
//
// Note that if there are multiple parameters with the same name, the parameters are
// retrieved and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequest(key string, def ...any) *gvar.Var {
	value := r.GetParam(key)
	if value.IsNil() {
		value = r.GetForm(key)
	}
	if value.IsNil() {
		r.parseBody()
		if len(r.bodyMap) > 0 {
			if v := r.bodyMap[key]; v != nil {
				value = gvar.New(v)
			}
		}
	}
	if value.IsNil() {
		value = r.GetQuery(key)
	}
	if value.IsNil() {
		value = r.GetRouter(key)
	}
	if !value.IsNil() {
		return value
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetRequestMap retrieves and returns all parameters passed from the client and custom params
// as the map, no matter what HTTP method the client is using. The parameter `kvMap` specifies
// the keys retrieving from client parameters, the associated values are the default values
// if the client does not pass the according keys.
//
// GetRequestMap is one of the most commonly used functions for retrieving parameters.
//
// Note that if there are multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequestMap(kvMap ...map[string]any) map[string]any {
	r.parseQuery()
	r.parseForm()
	r.parseBody()
	var (
		ok, filter bool
	)
	if len(kvMap) > 0 && kvMap[0] != nil {
		filter = true
	}
	m := make(map[string]any)
	for k, v := range r.routerMap {
		if filter {
			if _, ok = kvMap[0][k]; !ok {
				continue
			}
		}
		m[k] = v
	}
	for k, v := range r.queryMap {
		if filter {
			if _, ok = kvMap[0][k]; !ok {
				continue
			}
		}
		m[k] = v
	}
	for k, v := range r.formMap {
		if filter {
			if _, ok = kvMap[0][k]; !ok {
				continue
			}
		}
		m[k] = v
	}
	for k, v := range r.bodyMap {
		if filter {
			if _, ok = kvMap[0][k]; !ok {
				continue
			}
		}
		m[k] = v
	}
	for k, v := range r.paramsMap {
		if filter {
			if _, ok = kvMap[0][k]; !ok {
				continue
			}
		}
		m[k] = v
	}
	// File uploading.
	if r.MultipartForm != nil {
		for name := range r.MultipartForm.File {
			if uploadFiles := r.GetUploadFiles(name); len(uploadFiles) == 1 {
				m[name] = uploadFiles[0]
			} else {
				m[name] = uploadFiles
			}
		}
	}
	// Check none exist parameters and assign it with default value.
	if filter {
		for k, v := range kvMap[0] {
			if _, ok = m[k]; !ok {
				m[k] = v
			}
		}
	}
	return m
}

// GetRequestMapStrStr retrieve and returns all parameters passed from the client and custom
// params as map[string]string, no matter what HTTP method the client is using. The parameter
// `kvMap` specifies the keys retrieving from client parameters, the associated values are the
// default values if the client does not pass.
func (r *Request) GetRequestMapStrStr(kvMap ...map[string]any) map[string]string {
	requestMap := r.GetRequestMap(kvMap...)
	if len(requestMap) > 0 {
		m := make(map[string]string, len(requestMap))
		for k, v := range requestMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetRequestMapStrVar retrieve and returns all parameters passed from the client and custom
// params as map[string]*gvar.Var, no matter what HTTP method the client is using. The parameter
// `kvMap` specifies the keys retrieving from client parameters, the associated values are the
// default values if the client does not pass.
func (r *Request) GetRequestMapStrVar(kvMap ...map[string]any) map[string]*gvar.Var {
	requestMap := r.GetRequestMap(kvMap...)
	if len(requestMap) > 0 {
		m := make(map[string]*gvar.Var, len(requestMap))
		for k, v := range requestMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetRequestStruct retrieves all parameters passed from the client and custom params no matter
// what HTTP method the client is using, and converts them to give the struct object. Note that
// the parameter `pointer` is a pointer to the struct object.
// The optional parameter `mapping` is used to specify the key to attribute mapping.
func (r *Request) GetRequestStruct(pointer any, mapping ...map[string]string) error {
	_, err := r.doGetRequestStruct(pointer, mapping...)
	return err
}

func (r *Request) doGetRequestStruct(pointer any, mapping ...map[string]string) (data map[string]any, err error) {
	data = r.GetRequestMap()
	if err = r.GetError(); err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]any{}
	}
	// A field tagged with `in:"body"` represents the complete request body. Before converting,
	// remove the request parameter that is named exactly after the field tag name, as the tag
	// name is matched with a higher priority than the field name by the struct converting,
	// which would otherwise overwrite the body array. Both names are resolved at router
	// registering time, see checkAndCreateReqBodyField.
	if r.serveHandler.Handler.Info.ReqBodyFieldName != "" {
		delete(data, r.serveHandler.Handler.Info.ReqBodyFieldTagName)
	}

	// `in` Tag Struct values.
	if err = r.mergeInTagStructValue(data); err != nil {
		return data, nil
	}

	// Default struct values.
	if err = r.mergeDefaultStructValue(data, pointer); err != nil {
		return data, nil
	}

	// The request struct field tagged with `in:"body"` receives the whole request body, which
	// is a JSON array instead of being split into the request parameters.
	if bodyFieldName := r.serveHandler.Handler.Info.ReqBodyFieldName; bodyFieldName != "" {
		if r.bodyArray != nil {
			data[bodyFieldName] = r.bodyArray
		} else if r.bodyMap != nil || r.MultipartForm != nil {
			// The JSON object, the form parameters and the multipart forms are not acceptable
			// for such field, which are reported as an invalid parameter instead of being
			// silently ignored with the field left as a nil slice.
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`the request body should be a JSON array for the request struct field "%s" tagged with in:"body"`,
				bodyFieldName,
			)
		} else {
			// There's no body at all. The nil value occupies the field name, so that the field
			// is bound to its zero value before the request parameters of similar names (case
			// or symbol variants) could be fuzzy matched to it.
			data[bodyFieldName] = nil
		}
	}

	return data, gconv.Struct(data, pointer, mapping...)
}

// mergeDefaultStructValue merges the request parameters with default values from struct tag definition.
func (r *Request) mergeDefaultStructValue(data map[string]any, pointer any) error {
	fields := r.serveHandler.Handler.Info.ReqStructFields
	if len(fields) > 0 {
		for _, field := range fields {
			if tagValue := field.TagDefault(); tagValue != "" {
				mergeTagValueWithFoundKey(data, false, field.Name(), field.Name(), tagValue)
			}
		}
		return nil
	}

	// provide non strict routing
	tagFields, err := gstructs.TagFields(pointer, defaultValueTags)
	if err != nil {
		return err
	}
	if len(tagFields) > 0 {
		for _, field := range tagFields {
			mergeTagValueWithFoundKey(data, false, field.Name(), field.Name(), field.TagValue)
		}
	}

	return nil
}

// mergeInTagStructValue merges the request parameters with header or cookie values from struct `in` tag definition.
func (r *Request) mergeInTagStructValue(data map[string]any) error {
	fields := r.serveHandler.Handler.Info.ReqStructFields
	if len(fields) > 0 {
		var (
			headerMap = make(map[string]any)
			cookieMap = make(map[string]any)
		)

		for k, v := range r.Header {
			if len(v) > 0 {
				headerMap[k] = v[0]
			}
		}

		for _, cookie := range r.Cookies() {
			cookieMap[cookie.Name] = cookie.Value
		}

		for _, field := range fields {
			var (
				foundKey   string
				foundValue any
			)
			if tagValue := field.TagIn(); tagValue != "" {
				findKey := field.TagPriorityName()
				switch tagValue {
				case goai.ParameterInHeader:
					foundKey, foundValue = gutil.MapPossibleItemByKey(headerMap, findKey)
				case goai.ParameterInCookie:
					foundKey, foundValue = gutil.MapPossibleItemByKey(cookieMap, findKey)
				}
				if foundKey != "" {
					mergeTagValueWithFoundKey(data, true, foundKey, field.Name(), foundValue)
				}
			}
		}
	}
	return nil
}

// mergeTagValueWithFoundKey merges the request parameters when the key does not exist in the map or overwritten is true or the value is nil.
func mergeTagValueWithFoundKey(data map[string]any, overwritten bool, findKey string, fieldName string, tagValue any) {
	if foundKey, foundValue := gutil.MapPossibleItemByKey(data, findKey); foundKey == "" {
		data[fieldName] = tagValue
	} else {
		if overwritten || foundValue == nil {
			data[foundKey] = tagValue
		}
	}
}

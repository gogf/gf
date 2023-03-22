// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
)

// GetRequest retrieves and returns the parameter named `key` passed from the client and
// custom params as interface{}, no matter what HTTP method the client is using. The
// parameter `def` specifies the default value if the `key` does not exist.
//
// GetRequest is one of the most commonly used functions for retrieving parameters.
//
// Note that if there are multiple parameters with the same name, the parameters are
// retrieved and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequest(key string, def ...interface{}) *gvar.Var {
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
func (r *Request) GetRequestMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseQuery()
	r.parseForm()
	r.parseBody()
	var (
		ok, filter bool
	)
	if len(kvMap) > 0 && kvMap[0] != nil {
		filter = true
	}
	m := make(map[string]interface{})
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
func (r *Request) GetRequestMapStrStr(kvMap ...map[string]interface{}) map[string]string {
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
func (r *Request) GetRequestMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
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
func (r *Request) GetRequestStruct(pointer interface{}, mapping ...map[string]string) error {
	_, err := r.doGetRequestStruct(pointer, mapping...)
	return err
}

func (r *Request) doGetRequestStruct(pointer interface{}, mapping ...map[string]string) (data map[string]interface{}, err error) {
	data = r.GetRequestMap()
	if data == nil {
		data = map[string]interface{}{}
	}
	// Default struct values.
	if err = r.mergeDefaultStructValue(data, pointer); err != nil {
		return data, nil
	}
	// `in` Tag Struct values.
	if err = r.mergeInTagStructValue(data, pointer); err != nil {
		return data, nil
	}
	return data, gconv.Struct(data, pointer, mapping...)
}

// mergeDefaultStructValue merges the request parameters with default values from struct tag definition.
func (r *Request) mergeDefaultStructValue(data map[string]interface{}, pointer interface{}) error {
	// TODO: https://github.com/gogf/gf/pull/2450
	tagFields, err := gstructs.TagFields(pointer, defaultValueTags)
	if err != nil {
		return err
	}
	if len(tagFields) > 0 {
		var (
			foundKey   string
			foundValue interface{}
		)
		for _, field := range tagFields {
			foundKey, foundValue = gutil.MapPossibleItemByKey(data, field.Name())
			if foundKey == "" {
				data[field.Name()] = field.TagValue
			} else {
				if empty.IsEmpty(foundValue) {
					data[foundKey] = field.TagValue
				}
			}
		}
	}
	return nil
}

// mergeInTagStructValue merges the request parameters with header or cookie values from struct `in` tag definition.
func (r *Request) mergeInTagStructValue(data map[string]interface{}, pointer interface{}) error {
	// TODO: https://github.com/gogf/gf/pull/2450
	tagFields, err := gstructs.TagFields(pointer, []string{gtag.In})
	if err != nil {
		return err
	}
	if len(tagFields) > 0 {
		var (
			foundKey   string
			foundValue interface{}
			headerMap  = make(map[string]interface{})
			cookieMap  = make(map[string]interface{})
		)

		for k, v := range r.Header {
			if v != nil && len(v) > 0 {
				headerMap[k] = v[0]
			}
		}

		for _, cookie := range r.Cookies() {
			cookieMap[cookie.Name] = cookie.Value
		}

		for _, field := range tagFields {
			switch field.TagValue {
			case goai.ParameterInHeader:
				foundHeaderKey, foundHeaderValue := gutil.MapPossibleItemByKey(headerMap, field.Name())
				if foundHeaderKey != "" {
					foundKey, foundValue = gutil.MapPossibleItemByKey(data, foundHeaderKey)
					if foundKey == "" {
						data[field.Name()] = foundHeaderValue
					} else {
						if empty.IsEmpty(foundValue) {
							data[foundKey] = foundHeaderValue
						}
					}
				}
			case goai.ParameterInCookie:
				foundCookieKey, foundCookieValue := gutil.MapPossibleItemByKey(cookieMap, field.Name())
				if foundCookieKey != "" {
					foundKey, foundValue = gutil.MapPossibleItemByKey(data, foundCookieKey)
					if foundKey == "" {
						data[field.Name()] = foundCookieValue
					} else {
						if empty.IsEmpty(foundValue) {
							data[foundKey] = foundCookieValue
						}
					}
				}
			}
		}
	}
	return nil
}

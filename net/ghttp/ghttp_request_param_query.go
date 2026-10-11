// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// SetQuery sets custom query value with key-value pairs.
func (r *Request) SetQuery(key string, value any) {
	r.parseQuery()
	if r.queryMap == nil {
		r.queryMap = make(map[string]any)
	}
	r.queryMap[key] = value
}

// GetQuery retrieves and return parameter with the given name `key` from query string
// and request body. It returns `def` if `key` does not exist in the query and `def` is given,
// or else it returns nil.
//
// Note that if there are multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: query > body.
func (r *Request) GetQuery(key string, def ...any) *gvar.Var {
	r.parseQuery()
	if len(r.queryMap) > 0 {
		if value, ok := r.queryMap[key]; ok {
			return gvar.New(value)
		}
	}
	if r.Method == http.MethodGet {
		r.parseBody()
	}
	if len(r.bodyMap) > 0 {
		if v, ok := r.bodyMap[key]; ok {
			return gvar.New(v)
		}
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetQueryMap retrieves and returns all parameters passed from the client using HTTP GET method
// as the map. The parameter `kvMap` specifies the keys retrieving from client parameters,
// the associated values are the default values if the client does not pass.
//
// Note that if there are multiple parameters with the same name, the parameters are retrieved and overwrote
// in order of priority: query > body.
func (r *Request) GetQueryMap(kvMap ...map[string]any) map[string]any {
	r.parseQuery()
	if r.Method == http.MethodGet {
		r.parseBody()
	}
	var m map[string]any
	if len(kvMap) > 0 && kvMap[0] != nil {
		if len(r.queryMap) == 0 && len(r.bodyMap) == 0 {
			return kvMap[0]
		}
		m = make(map[string]any, len(kvMap[0]))
		if len(r.bodyMap) > 0 {
			for k, v := range kvMap[0] {
				if postValue, ok := r.bodyMap[k]; ok {
					m[k] = postValue
				} else {
					m[k] = v
				}
			}
		}
		if len(r.queryMap) > 0 {
			for k, v := range kvMap[0] {
				if postValue, ok := r.queryMap[k]; ok {
					m[k] = postValue
				} else {
					m[k] = v
				}
			}
		}
	} else {
		m = make(map[string]any, len(r.queryMap)+len(r.bodyMap))
		for k, v := range r.bodyMap {
			m[k] = v
		}
		for k, v := range r.queryMap {
			m[k] = v
		}
	}
	return m
}

// GetQueryMapStrStr retrieves and returns all parameters passed from the client using the HTTP GET method as a
//
//	map[string]string. The parameter `kvMap` specifies the keys
//
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrStr(kvMap ...map[string]any) map[string]string {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]string, len(queryMap))
		for k, v := range queryMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetQueryMapStrVar retrieves and returns all parameters passed from the client using the HTTP GET method
// as map[string]*gvar.Var. The parameter `kvMap` specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrVar(kvMap ...map[string]any) map[string]*gvar.Var {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]*gvar.Var, len(queryMap))
		for k, v := range queryMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetQueryStruct retrieves all parameters passed from the client using the HTTP GET method
// and converts them to a given struct object. Note that the parameter `pointer` is a pointer
// to the struct object. The optional parameter `mapping` is used to specify the key to
// attribute mapping.
func (r *Request) GetQueryStruct(pointer any, mapping ...map[string]string) error {
	_, err := r.doGetQueryStruct(pointer, mapping...)
	return err
}

func (r *Request) doGetQueryStruct(pointer any, mapping ...map[string]string) (data map[string]any, err error) {
	r.parseQuery()
	data = r.GetQueryMap()
	if err = r.GetError(); err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]any{}
	}
	r.mergeRepeatedQueryValues(data, pointer)
	if err = r.mergeDefaultStructValue(data, pointer); err != nil {
		return data, nil
	}
	return data, gconv.Struct(data, pointer, mapping...)
}

// mergeRepeatedQueryValues merges all the values of the repeated query parameters into `data`,
// like `names=a&names=b`, for the struct fields of slice type in `pointer`.
//
// The query string is parsed by gstr.Parse, which keeps only the last value of a repeated plain
// key, so the slice field would receive only that last value without this merging. The value in
// `data` is replaced only if it is still the last value from the query string, as it might be
// overridden by the parameters of higher priority, like the form, body or custom ones.
func (r *Request) mergeRepeatedQueryValues(data map[string]any, pointer any) {
	if r.URL.RawQuery == "" {
		return
	}
	repeatedValues := make(map[string][]string)
	for key, values := range r.URL.Query() {
		// The keys like `names[]` are already parsed as slice by gstr.Parse.
		if len(values) > 1 && !strings.Contains(key, "[") {
			repeatedValues[key] = values
		}
	}
	if len(repeatedValues) == 0 {
		return
	}
	fields, err := gstructs.Fields(gstructs.FieldsInput{
		Pointer:         pointer,
		RecursiveOption: gstructs.RecursiveOptionEmbedded,
	})
	if err != nil {
		// The `pointer` is not a struct, which has no field to merge and is left to the struct
		// converting as it is.
		return
	}
	for _, field := range fields {
		fieldType := field.Field.Type
		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}
		// The []byte field receives the value as a whole like the string field.
		if fieldType.Kind() != reflect.Slice || fieldType.Elem().Kind() == reflect.Uint8 {
			continue
		}
		foundKey, foundValue := gutil.MapPossibleItemByKey(data, field.TagPriorityName())
		values, ok := repeatedValues[foundKey]
		if !ok || foundValue != values[len(values)-1] {
			continue
		}
		data[foundKey] = values
	}
}

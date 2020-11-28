// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

// GetRequest retrieves and returns the parameter named <key> passed from client and
// custom params as interface{}, no matter what HTTP method the client is using. The
// parameter <def> specifies the default value if the <key> does not exist.
//
// GetRequest is one of the most commonly used functions for retrieving parameters.
//
// Note that if there're multiple parameters with the same name, the parameters are
// retrieved and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequest(key string, def ...interface{}) interface{} {
	value := r.GetParam(key)
	if value == nil {
		value = r.GetForm(key)
	}
	if value == nil {
		r.parseBody()
		if len(r.bodyMap) > 0 {
			value = r.bodyMap[key]
		}
	}
	if value == nil {
		value = r.GetQuery(key)
	}
	if value == nil {
		value = r.GetRouterValue(key)
	}
	if value != nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return value
}

// GetRequestVar retrieves and returns the parameter named <key> passed from client and
// custom params as gvar.Var, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRequest(key, def...))
}

// GetRequestString retrieves and returns the parameter named <key> passed from client and
// custom params as string, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestString(key string, def ...interface{}) string {
	return r.GetRequestVar(key, def...).String()
}

// GetRequestBool retrieves and returns the parameter named <key> passed from client and
// custom params as bool, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestBool(key string, def ...interface{}) bool {
	return r.GetRequestVar(key, def...).Bool()
}

// GetRequestInt retrieves and returns the parameter named <key> passed from client and
// custom params as int, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt(key string, def ...interface{}) int {
	return r.GetRequestVar(key, def...).Int()
}

// GetRequestInt32 retrieves and returns the parameter named <key> passed from client and
// custom params as int32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt32(key string, def ...interface{}) int32 {
	return r.GetRequestVar(key, def...).Int32()
}

// GetRequestInt64 retrieves and returns the parameter named <key> passed from client and
// custom params as int64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt64(key string, def ...interface{}) int64 {
	return r.GetRequestVar(key, def...).Int64()
}

// GetRequestInts retrieves and returns the parameter named <key> passed from client and
// custom params as []int, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInts(key string, def ...interface{}) []int {
	return r.GetRequestVar(key, def...).Ints()
}

// GetRequestUint retrieves and returns the parameter named <key> passed from client and
// custom params as uint, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint(key string, def ...interface{}) uint {
	return r.GetRequestVar(key, def...).Uint()
}

// GetRequestUint32 retrieves and returns the parameter named <key> passed from client and
// custom params as uint32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestVar(key, def...).Uint32()
}

// GetRequestUint64 retrieves and returns the parameter named <key> passed from client and
// custom params as uint64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestVar(key, def...).Uint64()
}

// GetRequestFloat32 retrieves and returns the parameter named <key> passed from client and
// custom params as float32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestVar(key, def...).Float32()
}

// GetRequestFloat64 retrieves and returns the parameter named <key> passed from client and
// custom params as float64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestVar(key, def...).Float64()
}

// GetRequestFloats retrieves and returns the parameter named <key> passed from client and
// custom params as []float64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestVar(key, def...).Floats()
}

// GetRequestArray retrieves and returns the parameter named <key> passed from client and
// custom params as []string, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestArray(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

// GetRequestStrings retrieves and returns the parameter named <key> passed from client and
// custom params as []string, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestStrings(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

// GetRequestInterfaces retrieves and returns the parameter named <key> passed from client
// and custom params as []interface{}, no matter what HTTP method the client is using. The
// parameter <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestVar(key, def...).Interfaces()
}

// GetRequestMap retrieves and returns all parameters passed from client and custom params
// as map, no matter what HTTP method the client is using. The parameter <kvMap> specifies
// the keys retrieving from client parameters, the associated values are the default values
// if the client does not pass the according keys.
//
// GetRequestMap is one of the most commonly used functions for retrieving parameters.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequestMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseQuery()
	r.parseForm()
	r.parseBody()
	var ok, filter bool
	var length int
	if len(kvMap) > 0 && kvMap[0] != nil {
		length = len(kvMap[0])
		filter = true
	} else {
		length = len(r.routerMap) + len(r.queryMap) + len(r.formMap) + len(r.bodyMap) + len(r.paramsMap)
	}
	m := make(map[string]interface{}, length)
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

// GetRequestMapStrStr retrieves and returns all parameters passed from client and custom
// params as map[string]string, no matter what HTTP method the client is using. The parameter
// <kvMap> specifies the keys retrieving from client parameters, the associated values are the
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

// GetRequestMapStrVar retrieves and returns all parameters passed from client and custom
// params as map[string]*gvar.Var, no matter what HTTP method the client is using. The parameter
// <kvMap> specifies the keys retrieving from client parameters, the associated values are the
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

// GetRequestStruct retrieves all parameters passed from client and custom params no matter
// what HTTP method the client is using, and converts them to given struct object. Note that
// the parameter <pointer> is a pointer to the struct object.
// The optional parameter <mapping> is used to specify the key to attribute mapping.
func (r *Request) GetRequestStruct(pointer interface{}, mapping ...map[string]string) error {
	data := r.GetRequestMap()
	if data == nil {
		data = map[string]interface{}{}
	}
	if err := r.mergeDefaultStructValue(data, pointer); err != nil {
		return nil
	}
	return gconv.Struct(data, pointer, mapping...)
}

// mergeDefaultStructValue merges the request parameters with default values from struct tag definition.
func (r *Request) mergeDefaultStructValue(data map[string]interface{}, pointer interface{}) error {
	tagFields, err := structs.TagFields(pointer, defaultValueTags)
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

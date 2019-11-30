// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
)

// GetRequestVar retrieves and returns the parameter named <key> passed from client as interface{},
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
//
// Note that the parameter is retrieved in order of: router->get/body->post/body->param.
func (r *Request) GetRequest(key string, def ...interface{}) interface{} {
	v := r.GetRouterValue(key)
	if v == nil {
		r.ParseQuery()
		if len(r.queryMap) > 0 {
			v, _ = r.queryMap[key]
		}
	}
	if v == nil {
		v = r.GetPost(key)
	}
	if v == nil {
		v = r.GetParam(key)
	}
	if v != nil {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return v
}

// GetRequestVar retrieves and returns the parameter named <key> passed from client as *gvar.Var,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRequest(key, def...))
}

// GetRequestString retrieves and returns the parameter named <key> passed from client as string,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestString(key string, def ...interface{}) string {
	return r.GetRequestVar(key, def...).String()
}

// GetRequestBool retrieves and returns the parameter named <key> passed from client as bool,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestBool(key string, def ...interface{}) bool {
	return r.GetRequestVar(key, def...).Bool()
}

// GetRequestInt retrieves and returns the parameter named <key> passed from client as int,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestInt(key string, def ...interface{}) int {
	return r.GetRequestVar(key, def...).Int()
}

// GetRequestInt32 retrieves and returns the parameter named <key> passed from client as int32,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestInt32(key string, def ...interface{}) int32 {
	return r.GetRequestVar(key, def...).Int32()
}

// GetRequestInt64 retrieves and returns the parameter named <key> passed from client as int64,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestInt64(key string, def ...interface{}) int64 {
	return r.GetRequestVar(key, def...).Int64()
}

// GetRequestInts retrieves and returns the parameter named <key> passed from client as []int,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestInts(key string, def ...interface{}) []int {
	return r.GetRequestVar(key, def...).Ints()
}

// GetRequestUint retrieves and returns the parameter named <key> passed from client as uint,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestUint(key string, def ...interface{}) uint {
	return r.GetRequestVar(key, def...).Uint()
}

// GetRequestUint32 retrieves and returns the parameter named <key> passed from client as uint32,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestVar(key, def...).Uint32()
}

// GetRequestUint64 retrieves and returns the parameter named <key> passed from client as uint64,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestVar(key, def...).Uint64()
}

// GetRequestFloat32 retrieves and returns the parameter named <key> passed from client as float32,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestVar(key, def...).Float32()
}

// GetRequestFloat64 retrieves and returns the parameter named <key> passed from client as float64,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestVar(key, def...).Float64()
}

// GetRequestFloats retrieves and returns the parameter named <key> passed from client as []float64,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestVar(key, def...).Floats()
}

// GetRequestArray retrieves and returns the parameter named <key> passed from client as []string,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestArray(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

// GetRequestStrings retrieves and returns the parameter named <key> passed from client as []string,
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestStrings(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

// GetRequestInterfaces retrieves and returns the parameter named <key> passed from client as []interface{},
// no matter what HTTP method the client is using. The parameter <def> specifies the default value
// if the <key> does not exist.
func (r *Request) GetRequestInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestVar(key, def...).Interfaces()
}

// GetRequestMap retrieves and returns all parameters passed from client as map,
// no matter what HTTP method the client is using. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetRequestMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.ParseQuery()
	r.ParseForm()
	r.ParseBody()
	m := make(map[string]interface{}, len(r.queryMap)+len(r.formMap)+len(r.bodyMap))
	for k, v := range r.queryMap {
		m[k] = v
	}
	for k, v := range r.formMap {
		m[k] = v
	}
	for k, v := range r.bodyMap {
		m[k] = v
	}
	if len(kvMap) > 0 && kvMap[0] != nil {
		var ok bool
		for k, _ := range m {
			if _, ok = kvMap[0][k]; !ok {
				delete(m, k)
			}
		}
	}
	return m
}

// GetRequestMapStrStr retrieves and returns all parameters passed from client as map[string]string,
// no matter what HTTP method the client is using. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetRequestMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	requestMap := r.GetRequestMap(kvMap...)
	if len(requestMap) > 0 {
		m := make(map[string]string)
		for k, v := range requestMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetRequestMapStrVar retrieves and returns all parameters passed from client as map[string]*gvar.Var,
// no matter what HTTP method the client is using. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetRequestMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	requestMap := r.GetRequestMap(kvMap...)
	if len(requestMap) > 0 {
		m := make(map[string]*gvar.Var)
		for k, v := range requestMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetRequestToStruct retrieves all parameters passed from client no matter what HTTP method the client is using,
// and converts them to given struct object. Note that the parameter <pointer> is a pointer to the struct object.
// The optional parameter <mapping> is used to specify the key to attribute mapping.
func (r *Request) GetRequestToStruct(pointer interface{}, mapping ...map[string]string) error {
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.GetRequestMap(), pointer, tagMap)
}

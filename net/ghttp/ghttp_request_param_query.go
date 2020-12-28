// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"

	"github.com/gogf/gf/util/gconv"
)

// SetQuery sets custom query value with key-value pair.
func (r *Request) SetQuery(key string, value interface{}) {
	r.parseQuery()
	if r.queryMap == nil {
		r.queryMap = make(map[string]interface{})
	}
	r.queryMap[key] = value
}

// GetQuery retrieves and returns parameter with given name <key> from query string
// and request body. It returns <def> if <key> does not exist in the query and <def> is given,
// or else it returns nil.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: query > body.
func (r *Request) GetQuery(key string, def ...interface{}) interface{} {
	r.parseQuery()
	if len(r.queryMap) > 0 {
		if v, ok := r.queryMap[key]; ok {
			return v
		}
	}
	if r.Method == "GET" {
		r.parseBody()
	}
	if len(r.bodyMap) > 0 {
		if v, ok := r.bodyMap[key]; ok {
			return v
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (r *Request) GetQueryVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetQuery(key, def...))
}

func (r *Request) GetQueryString(key string, def ...interface{}) string {
	return r.GetQueryVar(key, def...).String()
}

func (r *Request) GetQueryBool(key string, def ...interface{}) bool {
	return r.GetQueryVar(key, def...).Bool()
}

func (r *Request) GetQueryInt(key string, def ...interface{}) int {
	return r.GetQueryVar(key, def...).Int()
}

func (r *Request) GetQueryInt32(key string, def ...interface{}) int32 {
	return r.GetQueryVar(key, def...).Int32()
}

func (r *Request) GetQueryInt64(key string, def ...interface{}) int64 {
	return r.GetQueryVar(key, def...).Int64()
}

func (r *Request) GetQueryInts(key string, def ...interface{}) []int {
	return r.GetQueryVar(key, def...).Ints()
}

func (r *Request) GetQueryUint(key string, def ...interface{}) uint {
	return r.GetQueryVar(key, def...).Uint()
}

func (r *Request) GetQueryUint32(key string, def ...interface{}) uint32 {
	return r.GetQueryVar(key, def...).Uint32()
}

func (r *Request) GetQueryUint64(key string, def ...interface{}) uint64 {
	return r.GetQueryVar(key, def...).Uint64()
}

func (r *Request) GetQueryFloat32(key string, def ...interface{}) float32 {
	return r.GetQueryVar(key, def...).Float32()
}

func (r *Request) GetQueryFloat64(key string, def ...interface{}) float64 {
	return r.GetQueryVar(key, def...).Float64()
}

func (r *Request) GetQueryFloats(key string, def ...interface{}) []float64 {
	return r.GetQueryVar(key, def...).Floats()
}

func (r *Request) GetQueryArray(key string, def ...interface{}) []string {
	return r.GetQueryVar(key, def...).Strings()
}

func (r *Request) GetQueryStrings(key string, def ...interface{}) []string {
	return r.GetQueryVar(key, def...).Strings()
}

func (r *Request) GetQueryInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetQueryVar(key, def...).Interfaces()
}

// GetQueryMap retrieves and returns all parameters passed from client using HTTP GET method
// as map. The parameter <kvMap> specifies the keys retrieving from client parameters,
// the associated values are the default values if the client does not pass.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved and overwrote
// in order of priority: query > body.
func (r *Request) GetQueryMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseQuery()
	if r.Method == "GET" {
		r.parseBody()
	}
	var m map[string]interface{}
	if len(kvMap) > 0 && kvMap[0] != nil {
		if len(r.queryMap) == 0 && len(r.bodyMap) == 0 {
			return kvMap[0]
		}
		m = make(map[string]interface{}, len(kvMap[0]))
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
		m = make(map[string]interface{}, len(r.queryMap)+len(r.bodyMap))
		for k, v := range r.bodyMap {
			m[k] = v
		}
		for k, v := range r.queryMap {
			m[k] = v
		}
	}
	return m
}

// GetQueryMapStrStr retrieves and returns all parameters passed from client using HTTP GET method
// as map[string]string. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrStr(kvMap ...map[string]interface{}) map[string]string {
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

// GetQueryMapStrVar retrieves and returns all parameters passed from client using HTTP GET method
// as map[string]*gvar.Var. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
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

// GetQueryStruct retrieves all parameters passed from client using HTTP GET method
// and converts them to given struct object. Note that the parameter <pointer> is a pointer
// to the struct object. The optional parameter <mapping> is used to specify the key to
// attribute mapping.
func (r *Request) GetQueryStruct(pointer interface{}, mapping ...map[string]string) error {
	r.parseQuery()
	data := r.GetQueryMap()
	if data == nil {
		data = map[string]interface{}{}
	}
	if err := r.mergeDefaultStructValue(data, pointer); err != nil {
		return nil
	}
	return gconv.Struct(data, pointer, mapping...)
}

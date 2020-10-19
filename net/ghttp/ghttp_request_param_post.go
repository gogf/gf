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

// GetPost retrieves and returns parameter <key> from form and body.
// It returns <def> if <key> does not exist in neither form nor body.
// It returns nil if <def> is not passed.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: form > body.
//
// Deprecated.
func (r *Request) GetPost(key string, def ...interface{}) interface{} {
	r.parseForm()
	if len(r.formMap) > 0 {
		if v, ok := r.formMap[key]; ok {
			return v
		}
	}
	r.parseBody()
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

// Deprecated.
func (r *Request) GetPostVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetPost(key, def...))
}

// Deprecated.
func (r *Request) GetPostString(key string, def ...interface{}) string {
	return r.GetPostVar(key, def...).String()
}

// Deprecated.
func (r *Request) GetPostBool(key string, def ...interface{}) bool {
	return r.GetPostVar(key, def...).Bool()
}

// Deprecated.
func (r *Request) GetPostInt(key string, def ...interface{}) int {
	return r.GetPostVar(key, def...).Int()
}

// Deprecated.
func (r *Request) GetPostInt32(key string, def ...interface{}) int32 {
	return r.GetPostVar(key, def...).Int32()
}

// Deprecated.
func (r *Request) GetPostInt64(key string, def ...interface{}) int64 {
	return r.GetPostVar(key, def...).Int64()
}

// Deprecated.
func (r *Request) GetPostInts(key string, def ...interface{}) []int {
	return r.GetPostVar(key, def...).Ints()
}

// Deprecated.
func (r *Request) GetPostUint(key string, def ...interface{}) uint {
	return r.GetPostVar(key, def...).Uint()
}

// Deprecated.
func (r *Request) GetPostUint32(key string, def ...interface{}) uint32 {
	return r.GetPostVar(key, def...).Uint32()
}

// Deprecated.
func (r *Request) GetPostUint64(key string, def ...interface{}) uint64 {
	return r.GetPostVar(key, def...).Uint64()
}

// Deprecated.
func (r *Request) GetPostFloat32(key string, def ...interface{}) float32 {
	return r.GetPostVar(key, def...).Float32()
}

// Deprecated.
func (r *Request) GetPostFloat64(key string, def ...interface{}) float64 {
	return r.GetPostVar(key, def...).Float64()
}

// Deprecated.
func (r *Request) GetPostFloats(key string, def ...interface{}) []float64 {
	return r.GetPostVar(key, def...).Floats()
}

// Deprecated.
func (r *Request) GetPostArray(key string, def ...interface{}) []string {
	return r.GetPostVar(key, def...).Strings()
}

// Deprecated.
func (r *Request) GetPostStrings(key string, def ...interface{}) []string {
	return r.GetPostVar(key, def...).Strings()
}

// Deprecated.
func (r *Request) GetPostInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetPostVar(key, def...).Interfaces()
}

// GetPostMap retrieves and returns all parameters in the form and body passed from client
// as map. The parameter <kvMap> specifies the keys retrieving from client parameters,
// the associated values are the default values if the client does not pass.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved and overwrote
// in order of priority: form > body.
//
// Deprecated.
func (r *Request) GetPostMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseForm()
	r.parseBody()
	var ok, filter bool
	if len(kvMap) > 0 && kvMap[0] != nil {
		filter = true
	}
	m := make(map[string]interface{}, len(r.formMap)+len(r.bodyMap))
	for k, v := range r.bodyMap {
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

// GetPostMapStrStr retrieves and returns all parameters in the form and body passed from client
// as map[string]string. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
//
// Deprecated.
func (r *Request) GetPostMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	postMap := r.GetPostMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]string, len(postMap))
		for k, v := range postMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetPostMapStrVar retrieves and returns all parameters in the form and body passed from client
// as map[string]*gvar.Var. The parameter <kvMap> specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
//
// Deprecated.
func (r *Request) GetPostMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	postMap := r.GetPostMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]*gvar.Var, len(postMap))
		for k, v := range postMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetPostStruct retrieves all parameters in the form and body passed from client
// and converts them to given struct object. Note that the parameter <pointer> is a pointer
// to the struct object. The optional parameter <mapping> is used to specify the key to
// attribute mapping.
//
// Deprecated.
func (r *Request) GetPostStruct(pointer interface{}, mapping ...map[string]string) error {
	return gconv.Struct(r.GetPostMap(), pointer, mapping...)
}

// GetPostToStruct is alias of GetQueryStruct. See GetPostStruct.
//
// Deprecated.
func (r *Request) GetPostToStruct(pointer interface{}, mapping ...map[string]string) error {
	return r.GetPostStruct(pointer, mapping...)
}

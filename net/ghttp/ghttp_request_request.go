// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// 初始化RAW请求参数
func (r *Request) initRaw() {
	if !r.parsedRaw {
		r.parsedRaw = true
		if raw := r.GetRawString(); len(raw) > 0 {
			r.rawVarMap, _ = gstr.Parse(raw)
		}
	}
	if r.rawVarMap == nil {
		r.rawVarMap = make(map[string]interface{})
	}
}

// 获得router、post或者get提交的参数值，如果有同名参数，
// 那么按照 router->get->post->param->OtherHttpMethod 优先级进行覆盖。
// 注意获得参数值可能是字符串、数组、Map三种类型。
func (r *Request) GetRequest(key string, def ...interface{}) interface{} {
	v := r.GetRouterValue(key)
	if v == nil {
		v = r.GetQuery(key)
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
	r.initRaw()
	v = r.rawVarMap[key]
	if v == nil && len(def) > 0 {
		return def[0]
	}
	return v
}

func (r *Request) GetRequestVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRequest(key, def...))
}

func (r *Request) GetRequestString(key string, def ...interface{}) string {
	return r.GetRequestVar(key, def...).String()
}

func (r *Request) GetRequestBool(key string, def ...interface{}) bool {
	return r.GetRequestVar(key, def...).Bool()
}

func (r *Request) GetRequestInt(key string, def ...interface{}) int {
	return r.GetRequestVar(key, def...).Int()
}

func (r *Request) GetRequestInt32(key string, def ...interface{}) int32 {
	return r.GetRequestVar(key, def...).Int32()
}

func (r *Request) GetRequestInt64(key string, def ...interface{}) int64 {
	return r.GetRequestVar(key, def...).Int64()
}

func (r *Request) GetRequestInts(key string, def ...interface{}) []int {
	return r.GetRequestVar(key, def...).Ints()
}

func (r *Request) GetRequestUint(key string, def ...interface{}) uint {
	return r.GetRequestVar(key, def...).Uint()
}

func (r *Request) GetRequestUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestVar(key, def...).Uint32()
}

func (r *Request) GetRequestUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestVar(key, def...).Uint64()
}

func (r *Request) GetRequestFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestVar(key, def...).Float32()
}

func (r *Request) GetRequestFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestVar(key, def...).Float64()
}

func (r *Request) GetRequestFloats(key string, def ...interface{}) []float64 {
	return r.GetRequestVar(key, def...).Floats()
}

func (r *Request) GetRequestArray(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

func (r *Request) GetRequestStrings(key string, def ...interface{}) []string {
	return r.GetRequestVar(key, def...).Strings()
}

func (r *Request) GetRequestInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetRequestVar(key, def...).Interfaces()
}

func (r *Request) GetRequestMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.initRaw()
	m := r.rawVarMap
	if len(kvMap) > 0 {
		m = make(map[string]interface{})
		for k, defValue := range kvMap[0] {
			if rawValue, ok := r.rawVarMap[k]; ok {
				m[k] = rawValue
			} else {
				m[k] = defValue
			}
		}
	}
	if m == nil {
		m = make(map[string]interface{})
	}
	for k, v := range r.GetPostMap(kvMap...) {
		m[k] = v
	}
	for k, v := range r.GetQueryMap(kvMap...) {
		m[k] = v
	}
	return m
}

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

func (r *Request) GetRequestToStruct(pointer interface{}, mapping ...map[string]string) error {
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.GetRequestMap(), pointer, tagMap)
}

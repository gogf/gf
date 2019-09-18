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
}

// 获得router、post或者get提交的参数，如果有同名参数，那么按照router->get->post优先级进行覆盖
func (r *Request) GetRequest(key string, def ...interface{}) interface{} {
	v := r.GetRouterValue(key)
	if v == nil {
		v = r.GetQuery(key)
	}
	if v == nil {
		v = r.GetPost(key)
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

func (r *Request) GetRequestInts(key string, def ...interface{}) []int {
	return r.GetRequestVar(key, def...).Ints()
}

func (r *Request) GetRequestUint(key string, def ...interface{}) uint {
	return r.GetRequestVar(key, def...).Uint()
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

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetRequestArray获取特定字段内容
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

// 将所有的request参数映射到struct属性上，参数object应当为一个struct对象的指针, mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetRequestToStruct(pointer interface{}, mapping ...map[string]string) error {
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.GetRequestMap(), pointer, tagMap)
}

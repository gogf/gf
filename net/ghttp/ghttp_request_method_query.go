// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"

	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/container/gvar"

	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/util/gconv"
)

func (r *Request) initGet() {
	if !r.parsedGet {
		r.parsedGet = true
		if r.URL.RawQuery != "" {
			r.getMap, _ = gstr.Parse(r.URL.RawQuery)
		} else if strings.EqualFold(r.Method, "GET") {
			r.parsedRaw = true
			if raw := r.GetRawString(); len(raw) > 0 {
				r.getMap, _ = gstr.Parse(raw)
			}
		}
	}
	if r.getMap == nil {
		r.getMap = make(map[string]interface{})
	}
}

func (r *Request) SetQuery(key string, value interface{}) {
	r.initGet()
	r.getMap[key] = value
}

func (r *Request) GetQuery(key string, def ...interface{}) interface{} {
	r.initGet()
	if v, ok := r.getMap[key]; ok {
		return v
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

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值。
// 当不指定键值对关联数组时，默认获取GET方式提交的所有的提交键值对数据。
func (r *Request) GetQueryMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.initGet()
	if len(kvMap) > 0 {
		m := make(map[string]interface{})
		for k, defValue := range kvMap[0] {
			if queryValue, ok := r.getMap[k]; ok {
				m[k] = queryValue
			} else {
				m[k] = defValue
			}
		}
		return m
	} else {
		return r.getMap
	}
}

func (r *Request) GetQueryMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]string)
		for k, v := range queryMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetQueryMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]*gvar.Var)
		for k, v := range queryMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// 将所有的get参数映射到struct属性上，参数object应当为一个struct对象的指针, mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetQueryToStruct(pointer interface{}, mapping ...map[string]string) error {
	r.initGet()
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.getMap, pointer, tagMap)
}

// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"strings"
)

func (r *Request) initPut() {
	if !r.parsedPut {
		r.parsedPut = true
		if strings.EqualFold(r.Method, "PUT") {
			r.parsedRaw = true
			if raw := r.GetRawString(); len(raw) > 0 {
				r.putMap, _ = gstr.Parse(raw)
			}
		}
	}
	if r.putMap == nil {
		r.putMap = make(map[string]interface{})
	}
}

func (r *Request) SetPut(key string, value interface{}) {
	r.initPut()
	r.putMap[key] = value
}

func (r *Request) GetPut(key string, def ...interface{}) interface{} {
	r.initPut()
	if v, ok := r.putMap[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (r *Request) GetPutVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetPut(key, def...))
}

func (r *Request) GetPutString(key string, def ...interface{}) string {
	return r.GetPutVar(key, def...).String()
}

func (r *Request) GetPutBool(key string, def ...interface{}) bool {
	return r.GetPutVar(key, def...).Bool()
}

func (r *Request) GetPutInt(key string, def ...interface{}) int {
	return r.GetPutVar(key, def...).Int()
}

func (r *Request) GetPutInt32(key string, def ...interface{}) int32 {
	return r.GetPutVar(key, def...).Int32()
}

func (r *Request) GetPutInt64(key string, def ...interface{}) int64 {
	return r.GetPutVar(key, def...).Int64()
}

func (r *Request) GetPutInts(key string, def ...interface{}) []int {
	return r.GetPutVar(key, def...).Ints()
}

func (r *Request) GetPutUint(key string, def ...interface{}) uint {
	return r.GetPutVar(key, def...).Uint()
}

func (r *Request) GetPutUint32(key string, def ...interface{}) uint32 {
	return r.GetPutVar(key, def...).Uint32()
}

func (r *Request) GetPutUint64(key string, def ...interface{}) uint64 {
	return r.GetPutVar(key, def...).Uint64()
}

func (r *Request) GetPutFloat32(key string, def ...interface{}) float32 {
	return r.GetPutVar(key, def...).Float32()
}

func (r *Request) GetPutFloat64(key string, def ...interface{}) float64 {
	return r.GetPutVar(key, def...).Float64()
}

func (r *Request) GetPutFloats(key string, def ...interface{}) []float64 {
	return r.GetPutVar(key, def...).Floats()
}

func (r *Request) GetPutArray(key string, def ...interface{}) []string {
	return r.GetPutVar(key, def...).Strings()
}

func (r *Request) GetPutStrings(key string, def ...interface{}) []string {
	return r.GetPutVar(key, def...).Strings()
}

func (r *Request) GetPutInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetPutVar(key, def...).Interfaces()
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值。
// 当不指定键值对关联数组时，默认获取POST方式提交的所有的提交键值对数据。
func (r *Request) GetPutMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.initPut()
	if len(kvMap) > 0 {
		m := make(map[string]interface{})
		for k, defValue := range kvMap[0] {
			if putValue, ok := r.putMap[k]; ok {
				m[k] = putValue
			} else {
				m[k] = defValue
			}
		}
		return m
	} else {
		return r.putMap
	}
}

func (r *Request) GetPutMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	putMap := r.GetPutMap(kvMap...)
	if len(putMap) > 0 {
		m := make(map[string]string)
		for k, v := range putMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetPutMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	putMap := r.GetPutMap(kvMap...)
	if len(putMap) > 0 {
		m := make(map[string]*gvar.Var)
		for k, v := range putMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetPutToStruct(pointer interface{}, mapping ...map[string]string) error {
	r.initPut()
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.putMap, pointer, tagMap)
}

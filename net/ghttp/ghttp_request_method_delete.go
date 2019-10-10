// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

func (r *Request) initDelete() {
	if !r.parsedDelete {
		r.parsedDelete = true
		if strings.EqualFold(r.Method, "DELETE") {
			r.parsedRaw = true
			if raw := r.GetRawString(); len(raw) > 0 {
				r.deleteMap, _ = gstr.Parse(raw)
			}
		}
	}
	if r.deleteMap == nil {
		r.deleteMap = make(map[string]interface{})
	}
}

func (r *Request) SetDelete(key string, value interface{}) {
	r.initDelete()
	r.deleteMap[key] = value
}

func (r *Request) GetDelete(key string, def ...interface{}) interface{} {
	r.initDelete()
	if v, ok := r.deleteMap[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (r *Request) GetDeleteVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetDelete(key, def...))
}

func (r *Request) GetDeleteString(key string, def ...interface{}) string {
	return r.GetDeleteVar(key, def...).String()
}

func (r *Request) GetDeleteBool(key string, def ...interface{}) bool {
	return r.GetDeleteVar(key, def...).Bool()
}

func (r *Request) GetDeleteInt(key string, def ...interface{}) int {
	return r.GetDeleteVar(key, def...).Int()
}

func (r *Request) GetDeleteInt32(key string, def ...interface{}) int32 {
	return r.GetDeleteVar(key, def...).Int32()
}

func (r *Request) GetDeleteInt64(key string, def ...interface{}) int64 {
	return r.GetDeleteVar(key, def...).Int64()
}

func (r *Request) GetDeleteInts(key string, def ...interface{}) []int {
	return r.GetDeleteVar(key, def...).Ints()
}

func (r *Request) GetDeleteUint(key string, def ...interface{}) uint {
	return r.GetDeleteVar(key, def...).Uint()
}

func (r *Request) GetDeleteUint32(key string, def ...interface{}) uint32 {
	return r.GetDeleteVar(key, def...).Uint32()
}

func (r *Request) GetDeleteUint64(key string, def ...interface{}) uint64 {
	return r.GetDeleteVar(key, def...).Uint64()
}

func (r *Request) GetDeleteFloat32(key string, def ...interface{}) float32 {
	return r.GetDeleteVar(key, def...).Float32()
}

func (r *Request) GetDeleteFloat64(key string, def ...interface{}) float64 {
	return r.GetDeleteVar(key, def...).Float64()
}

func (r *Request) GetDeleteFloats(key string, def ...interface{}) []float64 {
	return r.GetDeleteVar(key, def...).Floats()
}

func (r *Request) GetDeleteArray(key string, def ...interface{}) []string {
	return r.GetDeleteVar(key, def...).Strings()
}

func (r *Request) GetDeleteStrings(key string, def ...interface{}) []string {
	return r.GetDeleteVar(key, def...).Strings()
}

func (r *Request) GetDeleteInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetDeleteVar(key, def...).Interfaces()
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值。
// 当不指定键值对关联数组时，默认获取POST方式提交的所有的提交键值对数据。
func (r *Request) GetDeleteMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.initDelete()
	if len(kvMap) > 0 {
		m := make(map[string]interface{})
		for k, defValue := range kvMap[0] {
			if deleteValue, ok := r.deleteMap[k]; ok {
				m[k] = deleteValue
			} else {
				m[k] = defValue
			}
		}
		return m
	} else {
		return r.deleteMap
	}
}

func (r *Request) GetDeleteMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	deleteMap := r.GetDeleteMap(kvMap...)
	if len(deleteMap) > 0 {
		m := make(map[string]string)
		for k, v := range deleteMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetDeleteMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	deleteMap := r.GetDeleteMap(kvMap...)
	if len(deleteMap) > 0 {
		m := make(map[string]*gvar.Var)
		for k, v := range deleteMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetDeleteToStruct(pointer interface{}, mapping ...map[string]string) error {
	r.initDelete()
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.deleteMap, pointer, tagMap)
}

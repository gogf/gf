// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"strings"

	"github.com/gogf/gf/encoding/gurl"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/structs"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

func (r *Request) initPost() {
	if !r.parsedPost {
		r.parsedPost = true
		if v := r.Header.Get("Content-Type"); v != "" && gstr.Contains(v, "multipart/") {
			// multipart/form-data, multipart/mixed
			r.ParseMultipartForm(r.Server.config.FormParsingMemory)
			if len(r.PostForm) > 0 {
				// 重新组织数据格式，使用统一的数据Parse方式
				params := ""
				for name, values := range r.PostForm {
					if len(values) == 1 {
						if len(params) > 0 {
							params += "&"
						}
						params += name + "=" + gurl.Encode(values[0])
					} else {
						if len(name) > 2 && name[len(name)-2:] == "[]" {
							name = name[:len(name)-2]
							for _, v := range values {
								if len(params) > 0 {
									params += "&"
								}
								params += name + "[]=" + gurl.Encode(v)
							}
						} else {
							if len(params) > 0 {
								params += "&"
							}
							params += name + "=" + gurl.Encode(values[len(values)-1])
						}
					}
				}
				r.postMap, _ = gstr.Parse(params)
			}
		} else if strings.EqualFold(r.Method, "POST") {
			r.parsedRaw = true
			if raw := r.GetRawString(); len(raw) > 0 {
				r.postMap, _ = gstr.Parse(raw)
			}
		}
	}
	if r.postMap == nil {
		r.postMap = make(map[string]interface{})
	}
}

func (r *Request) SetPost(key string, value interface{}) {
	r.initPost()
	r.postMap[key] = value
}

func (r *Request) GetPost(key string, def ...interface{}) interface{} {
	r.initPost()
	if v, ok := r.postMap[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

func (r *Request) GetPostVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetPost(key, def...))
}

func (r *Request) GetPostString(key string, def ...interface{}) string {
	return r.GetPostVar(key, def...).String()
}

func (r *Request) GetPostBool(key string, def ...interface{}) bool {
	return r.GetPostVar(key, def...).Bool()
}

func (r *Request) GetPostInt(key string, def ...interface{}) int {
	return r.GetPostVar(key, def...).Int()
}

func (r *Request) GetPostInt32(key string, def ...interface{}) int32 {
	return r.GetPostVar(key, def...).Int32()
}

func (r *Request) GetPostInt64(key string, def ...interface{}) int64 {
	return r.GetPostVar(key, def...).Int64()
}

func (r *Request) GetPostInts(key string, def ...interface{}) []int {
	return r.GetPostVar(key, def...).Ints()
}

func (r *Request) GetPostUint(key string, def ...interface{}) uint {
	return r.GetPostVar(key, def...).Uint()
}

func (r *Request) GetPostUint32(key string, def ...interface{}) uint32 {
	return r.GetPostVar(key, def...).Uint32()
}

func (r *Request) GetPostUint64(key string, def ...interface{}) uint64 {
	return r.GetPostVar(key, def...).Uint64()
}

func (r *Request) GetPostFloat32(key string, def ...interface{}) float32 {
	return r.GetPostVar(key, def...).Float32()
}

func (r *Request) GetPostFloat64(key string, def ...interface{}) float64 {
	return r.GetPostVar(key, def...).Float64()
}

func (r *Request) GetPostFloats(key string, def ...interface{}) []float64 {
	return r.GetPostVar(key, def...).Floats()
}

func (r *Request) GetPostArray(key string, def ...interface{}) []string {
	return r.GetPostVar(key, def...).Strings()
}

func (r *Request) GetPostStrings(key string, def ...interface{}) []string {
	return r.GetPostVar(key, def...).Strings()
}

func (r *Request) GetPostInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetPostVar(key, def...).Interfaces()
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值。
// 当不指定键值对关联数组时，默认获取POST方式提交的所有的提交键值对数据。
func (r *Request) GetPostMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.initPost()
	if len(kvMap) > 0 {
		m := make(map[string]interface{})
		for k, defValue := range kvMap[0] {
			if postValue, ok := r.postMap[k]; ok {
				m[k] = postValue
			} else {
				m[k] = defValue
			}
		}
		return m
	} else {
		return r.postMap
	}
}

func (r *Request) GetPostMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	postMap := r.GetPostMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]string)
		for k, v := range postMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetPostMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	postMap := r.GetPostMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]*gvar.Var)
		for k, v := range postMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

func (r *Request) GetPostToStruct(pointer interface{}, mapping ...map[string]string) error {
	r.initPost()
	tagMap := structs.TagMapName(pointer, paramTagPriority, true)
	if len(mapping) > 0 {
		for k, v := range mapping[0] {
			tagMap[k] = v
		}
	}
	return gconv.StructDeep(r.postMap, pointer, tagMap)
}

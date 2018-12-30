// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/util/gconv"
    "strings"
)

// 初始化GET请求参数
func (r *Request) initGet() {
    if !r.parsedGet {
        r.queryVars = r.URL.Query()
        if strings.EqualFold(r.Method, "GET") {
            if raw := r.GetRaw(); len(raw) > 0 {
                for _, item := range strings.Split(string(raw), "&") {
                    array                := strings.Split(item, "=")
                    r.queryVars[array[0]] = append(r.queryVars[array[0]], array[1])
                }
            }
        }
        r.parsedGet = true
    }
}

// 设置GET参数，仅在ghttp.Server内有效，**注意并发安全性**
func (r *Request) SetQuery(key string, value string) {
    r.initGet()
    r.queryVars[key] = []string{value}
}

// 添加GET参数，构成[]string
func (r *Request) AddQuery(key string, value string) {
    r.initGet()
    r.queryVars[key] = append(r.queryVars[key], value)
}

// 获得指定名称的get参数列表
func (r *Request) GetQuery(key string, def ... []string) []string {
    r.initGet()
    if v, ok := r.queryVars[key]; ok {
        return v
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

func (r *Request) GetQueryString(key string, def ... string) string {
    value := r.GetQuery(key)
    if value != nil && value[0] != "" {
        return value[0]
    }
    if len(def) > 0 {
        return def[0]
    }
    return ""
}

func (r *Request) GetQueryBool(key string, def ... bool) bool {
    value := r.GetQueryString(key)
    if value != "" {
        return gconv.Bool(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return false
}

func (r *Request) GetQueryInt(key string, def ... int) int {
    value := r.GetQueryString(key)
    if value != "" {
        return gconv.Int(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetQueryInts(key string, def ... []int) []int {
    value := r.GetQuery(key)
    if value != nil {
        return gconv.Ints(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

func (r *Request) GetQueryUint(key string, def ... uint) uint {
    value := r.GetQueryString(key)
    if value != "" {
        return gconv.Uint(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetQueryFloat32(key string, def ... float32) float32 {
    value := r.GetQueryString(key)
    if value != "" {
        return gconv.Float32(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetQueryFloat64(key string, def ... float64) float64 {
    value := r.GetQueryString(key)
    if value != "" {
        return gconv.Float64(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetQueryFloats(key string, def ... []float64) []float64 {
    value := r.GetQuery(key)
    if value != nil {
        return gconv.Floats(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

func (r *Request) GetQueryArray(key string, def ... []string) []string {
    return r.GetQuery(key, def...)
}

func (r *Request) GetQueryStrings(key string, def ... []string) []string {
    return r.GetQuery(key, def...)
}

func (r *Request) GetQueryInterfaces(key string, def ... []interface{}) []interface{} {
    value := r.GetQuery(key)
    if value != nil {
        return gconv.Interfaces(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetQueryMap(def ... map[string]string) map[string]string {
    r.initGet()
    m := make(map[string]string)
    for k, v := range r.queryVars {
        m[k] = v[0]
    }
    if len(def) > 0 {
        for k, v := range def[0] {
            if _, ok := m[k]; !ok {
                m[k] = v
            }
        }
    }
    return m
}

// 将所有的get参数映射到struct属性上，参数object应当为一个struct对象的指针, mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetQueryToStruct(object interface{}, mapping...map[string]string) {
    tagmap := r.getStructParamsTagMap(object)
    if len(mapping) > 0 {
        for k, v := range mapping[0] {
            tagmap[k] = v
        }
    }
    params := make(map[string]interface{})
    for k, v := range r.GetQueryMap() {
        params[k] = v
    }
    gconv.Struct(params, object, tagmap)
}
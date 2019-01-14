// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gvar"
)

// 获得router、post或者get提交的参数，如果有同名参数，那么按照router->get->post优先级进行覆盖
func (r *Request) GetRequest(key string, def ... []string) []string {
    v := r.GetRouterArray(key)
    if v == nil {
        v = r.GetQuery(key)
    }
    if v == nil {
        v = r.GetPost(key)
    }
    if v == nil && len(def) > 0 {
        return def[0]
    }
    return v
}

func (r *Request) GetRequestVar(key string, def ... interface{}) gvar.VarRead {
    value := r.GetRequest(key)
    if value != nil {
        return gvar.New(value[0], true)
    }
    if len(def) > 0 {
        return gvar.New(def[0], true)
    }
    return gvar.New(nil, true)
}

func (r *Request) GetRequestString(key string, def ... string) string {
    value := r.GetRequest(key)
    if value != nil && value[0] != "" {
        return value[0]
    }
    if len(def) > 0 {
        return def[0]
    }
    return ""
}

func (r *Request) GetRequestBool(key string, def ... bool) bool {
    value := r.GetRequestString(key)
    if value != "" {
        return gconv.Bool(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return false
}

func (r *Request) GetRequestInt(key string, def ... int) int {
    value := r.GetRequestString(key)
    if value != "" {
        return gconv.Int(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetRequestInts(key string, def ... []int) []int {
    value := r.GetRequest(key)
    if value != nil {
        return gconv.Ints(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

func (r *Request) GetRequestUint(key string, def ... uint) uint {
    value := r.GetRequestString(key)
    if value != "" {
        return gconv.Uint(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetRequestFloat32(key string, def ... float32) float32 {
    value := r.GetRequestString(key)
    if value != "" {
        return gconv.Float32(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetRequestFloat64(key string, def ... float64) float64 {
    value := r.GetRequestString(key)
    if value != "" {
        return gconv.Float64(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return 0
}

func (r *Request) GetRequestFloats(key string, def ... []float64) []float64 {
    value := r.GetRequest(key)
    if value != nil {
        return gconv.Floats(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

func (r *Request) GetRequestArray(key string, def ... []string) []string {
    return r.GetRequest(key, def...)
}

func (r *Request) GetRequestStrings(key string, def ... []string) []string {
    return r.GetRequest(key, def...)
}

func (r *Request) GetRequestInterfaces(key string, def ... []interface{}) []interface{} {
    value := r.GetRequest(key)
    if value != nil {
        return gconv.Interfaces(value)
    }
    if len(def) > 0 {
        return def[0]
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetRequestArray获取特定字段内容
func (r *Request) GetRequestMap(def...map[string]string) map[string]string {
    m := r.GetQueryMap()
    if len(m) == 0 {
        m = r.GetPostMap()
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

// 将所有的request参数映射到struct属性上，参数object应当为一个struct对象的指针, mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetRequestToStruct(object interface{}, mapping...map[string]string) {
    tagmap := r.getStructParamsTagMap(object)
    if len(mapping) > 0 {
        for k, v := range mapping[0] {
            tagmap[k] = v
        }
    }
    params := make(map[string]interface{})
    for k, v := range r.GetRequestMap() {
        params[k] = v
    }
    if len(params) == 0 {
        if j := r.GetJson(); j != nil {
            params = j.ToMap()
        }
    }
    gconv.Struct(params, object, tagmap)
}


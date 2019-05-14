// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/container/gvar"
)

// 获得router、post或者get提交的参数，如果有同名参数，那么按照router->get->post优先级进行覆盖
func (r *Request) GetRequest(key string, def...interface{}) []string {
    v := r.GetRouterArray(key)
    if v == nil {
        v = r.GetQuery(key)
    }
    if v == nil {
        v = r.GetPost(key)
    }
    if v == nil && len(def) > 0 {
        return gconv.Strings(def[0])
    }
    return v
}

func (r *Request) GetRequestVar(key string, def...interface{}) *gvar.Var {
    value := r.GetRequest(key, def...)
    if value != nil {
        return gvar.New(value[0], true)
    }
    return gvar.New(nil, true)
}

func (r *Request) GetRequestString(key string, def...interface{}) string {
    value := r.GetRequest(key, def...)
    if value != nil && value[0] != "" {
        return value[0]
    }
    return ""
}

func (r *Request) GetRequestBool(key string, def...interface{}) bool {
    value := r.GetRequestString(key, def...)
    if value != "" {
        return gconv.Bool(value)
    }
    return false
}

func (r *Request) GetRequestInt(key string, def...interface{}) int {
    value := r.GetRequestString(key, def...)
    if value != "" {
        return gconv.Int(value)
    }
    return 0
}

func (r *Request) GetRequestInts(key string, def...interface{}) []int {
    value := r.GetRequest(key, def...)
    if value != nil {
        return gconv.Ints(value)
    }
    return nil
}

func (r *Request) GetRequestUint(key string, def...interface{}) uint {
    value := r.GetRequestString(key, def...)
    if value != "" {
        return gconv.Uint(value)
    }
    return 0
}

func (r *Request) GetRequestFloat32(key string, def...interface{}) float32 {
    value := r.GetRequestString(key, def...)
    if value != "" {
        return gconv.Float32(value)
    }
    return 0
}

func (r *Request) GetRequestFloat64(key string, def...interface{}) float64 {
    value := r.GetRequestString(key, def...)
    if value != "" {
        return gconv.Float64(value)
    }
    return 0
}

func (r *Request) GetRequestFloats(key string, def...interface{}) []float64 {
    value := r.GetRequest(key, def...)
    if value != nil {
        return gconv.Floats(value)
    }
    return nil
}

func (r *Request) GetRequestArray(key string, def...interface{}) []string {
    return r.GetRequest(key, def...)
}

func (r *Request) GetRequestStrings(key string, def...interface{}) []string {
    return r.GetRequest(key, def...)
}

func (r *Request) GetRequestInterfaces(key string, def...interface{}) []interface{} {
    value := r.GetRequest(key, def...)
    if value != nil {
        return gconv.Interfaces(value)
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
func (r *Request) GetRequestToStruct(pointer interface{}, mapping...map[string]string) error {
    tagmap := r.getStructParamsTagMap(pointer)
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
    return gconv.Struct(params, pointer, tagmap)
}


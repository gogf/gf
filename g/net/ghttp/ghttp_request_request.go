// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/util/gconv"
)

// 获得router、post或者get提交的参数，如果有同名参数，那么按照router->get->post优先级进行覆盖
func (r *Request) GetRequest(k string) []string {
    v := r.GetRouterArray(k)
    if v == nil {
        v = r.GetQuery(k)
    }
    if v == nil {
        v = r.GetPost(k)
    }
    return v
}

func (r *Request) GetRequestString(k string) string {
    v := r.GetRequest(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetRequestBool(k string) bool {
    return gconv.Bool(r.GetRequestString(k))
}

func (r *Request) GetRequestInt(k string) int {
    return gconv.Int(r.GetRequestString(k))
}

func (r *Request) GetRequestUint(k string) uint {
    return gconv.Uint(r.GetRequestString(k))
}

func (r *Request) GetRequestFloat32(k string) float32 {
    return gconv.Float32(r.GetRequestString(k))
}

func (r *Request) GetRequestFloat64(k string) float64 {
    return gconv.Float64(r.GetRequestString(k))
}

func (r *Request) GetRequestArray(k string) []string {
    return r.GetRequest(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetRequestArray获取特定字段内容
func (r *Request) GetRequestMap(defaultMap...map[string]string) map[string]string {
    m := r.GetQueryMap()
    if len(defaultMap) == 0 {
        for k, v := range r.GetPostMap() {
            if _, ok := m[k]; !ok {
                m[k] = v
            }
        }
    } else {
        for k, v := range defaultMap[0] {
            v2 := r.GetRequest(k)
            if v2 != nil {
                m[k] = v2[0]
            } else {
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
    gconv.MapToStruct(params, object, tagmap)
}


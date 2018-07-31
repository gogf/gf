// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/util/gconv"
)

// 初始化GET请求参数
func (r *Request) initGet() {
    if !r.parsedGet.Val() {
        if len(r.queryVars) == 0 {
            r.queryVars = r.URL.Query()
        } else {
            for k, v := range r.URL.Query() {
                r.queryVars[k] = v
            }
        }
    }
}

// 设置GET参数，仅在ghttp.Server内有效，**注意并发安全性**
func (r *Request) SetQuery(k string, v string) {
    r.queryVars[k] = []string{v}
}

// 添加GET参数，构成[]string
func (r *Request) AddQuery(k string, v string) {
    r.queryVars[k] = append(r.queryVars[k], v)
}

// 获得指定名称的get参数列表
func (r *Request) GetQuery(k string) []string {
    r.initGet()
    if v, ok := r.queryVars[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetQueryBool(k string) bool {
    return gconv.Bool(r.Get(k))
}

func (r *Request) GetQueryInt(k string) int {
    return gconv.Int(r.Get(k))
}

func (r *Request) GetQueryUint(k string) uint {
    return gconv.Uint(r.Get(k))
}

func (r *Request) GetQueryFloat32(k string) float32 {
    return gconv.Float32(r.Get(k))
}

func (r *Request) GetQueryFloat64(k string) float64 {
    return gconv.Float64(r.Get(k))
}

func (r *Request) GetQueryString(k string) string {
    v := r.GetQuery(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetQueryArray(k string) []string {
    return r.GetQuery(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *Request) GetQueryMap(defaultMap...map[string]string) map[string]string {
    r.initGet()
    m := make(map[string]string)
    if len(defaultMap) == 0 {
        for k, v := range r.queryVars {
            m[k] = v[0]
        }
    } else {
        for k, v := range defaultMap[0] {
            v2 := r.GetQueryArray(k)
            if v2 == nil {
                m[k] = v
            } else {
                m[k] = v2[0]
            }
        }
    }
    return m
}


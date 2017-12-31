// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
package ghttp

import (
    "io/ioutil"
    "net/http"
    "net/url"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
)

// 请求对象
type ClientRequest struct {
    http.Request
    getvals  *url.Values    // GET参数
    Id       uint64         // 请求id(唯一)
    Cookie   *Cookie        // 与当前请求绑定的Cookie对象(并发安全)
    Session  *Session       // 与当前请求绑定的Session对象(并发安全)
}

// 获得指定名称的get参数列表
func (r *ClientRequest) GetQuery(k string) []string {
    if r.getvals == nil {
        values     := r.URL.Query()
        r.getvals = &values
    }
    if v, ok := (*r.getvals)[k]; ok {
        return v
    }
    return nil
}

func (r *ClientRequest) GetQueryBool(k string) bool {
    return gconv.Bool(r.GetQueryString(k))
}

func (r *ClientRequest) GetQueryInt(k string) int {
    return gconv.Int(r.GetQueryString(k))
}

func (r *ClientRequest) GetQueryUint(k string) uint {
    return gconv.Uint(r.GetQueryString(k))
}

func (r *ClientRequest) GetQueryFloat32(k string) float32 {
    return gconv.Float32(r.GetQueryString(k))
}

func (r *ClientRequest) GetQueryFloat64(k string) float64 {
    return gconv.Float64(r.GetQueryString(k))
}

func (r *ClientRequest) GetQueryString(k string) string {
    v := r.GetQuery(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetQueryArray(k string) []string {
    return r.GetQuery(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
func (r *ClientRequest) GetQueryMap(defaultMap map[string]string) map[string]string {
    m := make(map[string]string)
    for k, v := range defaultMap {
        v2 := r.GetQueryArray(k)
        if v2 == nil {
            m[k] = v
        } else {
            m[k] = v2[0]
        }
    }
    return m
}

// 获得post参数
func (r *ClientRequest) GetPost(k string) []string {
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
}

func (r *ClientRequest) GetPostBool(k string) bool {
    return gconv.Bool(r.GetPostString(k))
}

func (r *ClientRequest) GetPostInt(k string) int {
    return gconv.Int(r.GetPostString(k))
}

func (r *ClientRequest) GetPostUint(k string) uint {
    return gconv.Uint(r.GetPostString(k))
}

func (r *ClientRequest) GetPostFloat32(k string) float32 {
    return gconv.Float32(r.GetPostString(k))
}

func (r *ClientRequest) GetPostFloat64(k string) float64 {
    return gconv.Float64(r.GetPostString(k))
}

func (r *ClientRequest) GetPostString(k string) string {
    v := r.GetPost(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetPostArray(k string) []string {
    return r.GetPost(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetPostArray获取特定字段内容
func (r *ClientRequest) GetPostMap(defaultMap map[string]string) map[string]string {
    m := make(map[string]string)
    for k, v := range defaultMap {
        if v2, ok := r.PostForm[k]; ok {
            m[k] = v2[0]
        } else {
            m[k] = v
        }
    }
    return m
}

// 获得post或者get提交的参数，如果有同名参数，那么按照get->post优先级进行覆盖
func (r *ClientRequest) GetRequest(k string) []string {
    v := r.GetQuery(k)
    if v == nil {
        return r.GetPost(k)
    }
    return v
}

func (r *ClientRequest) GetRequestString(k string) string {
    v := r.GetRequest(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *ClientRequest) GetRequestBool(k string) bool {
    return gconv.Bool(r.GetRequestString(k))
}

func (r *ClientRequest) GetRequestInt(k string) int {
    return gconv.Int(r.GetRequestString(k))
}

func (r *ClientRequest) GetRequestUint(k string) uint {
    return gconv.Uint(r.GetRequestString(k))
}

func (r *ClientRequest) GetRequestFloat32(k string) float32 {
    return gconv.Float32(r.GetRequestString(k))
}

func (r *ClientRequest) GetRequestFloat64(k string) float64 {
    return gconv.Float64(r.GetRequestString(k))
}

func (r *ClientRequest) GetRequestArray(k string) []string {
    return r.GetRequest(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetRequestArray获取特定字段内容
func (r *ClientRequest) GetRequestMap(defaultMap map[string]string) map[string]string {
    m := make(map[string]string)
    for k, v := range defaultMap {
        v2 := r.GetRequest(k)
        if v2 != nil {
            m[k] = v2[0]
        } else {
            m[k] = v
        }
    }
    return m
}

// 获取原始请求输入字符串
func (r *ClientRequest) GetRaw() []byte {
    result, _ := ioutil.ReadAll(r.Body)
    return result
}

// 获取原始json请求输入字符串，并解析为json对象
func (r *ClientRequest) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != nil {
        if j, err := gjson.DecodeToJson(data); err == nil {
            return j
        }
    }
    return nil
}




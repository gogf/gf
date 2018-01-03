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
type Request struct {
    http.Request
    getvals  *url.Values     // GET参数
    Id       int             // 请求id(唯一)
    Server   *Server         // 请求关联的服务器对象
    Cookie   *Cookie         // 与当前请求绑定的Cookie对象(并发安全)
    Session  *Session        // 与当前请求绑定的Session对象(并发安全)
    Response *Response       // 对应请求的返回数据操作对象
}

// 获得指定名称的get参数列表
func (r *Request) GetQuery(k string) []string {
    if r.getvals == nil {
        values   := r.URL.Query()
        r.getvals = &values
    }
    if v, ok := (*r.getvals)[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetQueryBool(k string) bool {
    return gconv.Bool(r.GetQueryString(k))
}

func (r *Request) GetQueryInt(k string) int {
    return gconv.Int(r.GetQueryString(k))
}

func (r *Request) GetQueryUint(k string) uint {
    return gconv.Uint(r.GetQueryString(k))
}

func (r *Request) GetQueryFloat32(k string) float32 {
    return gconv.Float32(r.GetQueryString(k))
}

func (r *Request) GetQueryFloat64(k string) float64 {
    return gconv.Float64(r.GetQueryString(k))
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
func (r *Request) GetQueryMap(defaultMap map[string]string) map[string]string {
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
func (r *Request) GetPost(k string) []string {
    if v, ok := r.PostForm[k]; ok {
        return v
    }
    return nil
}

func (r *Request) GetPostBool(k string) bool {
    return gconv.Bool(r.GetPostString(k))
}

func (r *Request) GetPostInt(k string) int {
    return gconv.Int(r.GetPostString(k))
}

func (r *Request) GetPostUint(k string) uint {
    return gconv.Uint(r.GetPostString(k))
}

func (r *Request) GetPostFloat32(k string) float32 {
    return gconv.Float32(r.GetPostString(k))
}

func (r *Request) GetPostFloat64(k string) float64 {
    return gconv.Float64(r.GetPostString(k))
}

func (r *Request) GetPostString(k string) string {
    v := r.GetPost(k)
    if v == nil {
        return ""
    } else {
        return v[0]
    }
}

func (r *Request) GetPostArray(k string) []string {
    return r.GetPost(k)
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetPostArray获取特定字段内容
func (r *Request) GetPostMap(defaultMap map[string]string) map[string]string {
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
func (r *Request) GetRequest(k string) []string {
    v := r.GetQuery(k)
    if v == nil {
        return r.GetPost(k)
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
func (r *Request) GetRequestMap(defaultMap map[string]string) map[string]string {
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
func (r *Request) GetRaw() []byte {
    result, _ := ioutil.ReadAll(r.Body)
    return result
}

// 获取原始json请求输入字符串，并解析为json对象
func (r *Request) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != nil {
        if j, err := gjson.DecodeToJson(data); err == nil {
            return j
        }
    }
    return nil
}




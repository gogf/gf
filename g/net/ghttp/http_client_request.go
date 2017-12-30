// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
package ghttp

import (
    "io/ioutil"
    "gitee.com/johng/gf/g/encoding/gjson"
    "strconv"
)

// 获取当前请求的id
func (r *ClientRequest) Id() uint64 {
    return r.id
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

// 获取指定名称的参数int类型
func (r *ClientRequest) GetQueryInt(k string) int {
    v := r.GetQuery(k)
    if v == nil {
        return -1
    } else {
        if i, err := strconv.Atoi(v[0]); err != nil {
            return -1
        } else {
            return i
        }
    }
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
    v := r.GetQuery(k)
    if v == nil {
        return nil
    } else {
        return v
    }
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

func (r *ClientRequest) GetPostInt(k string) int {
    v := r.GetPost(k)
    if v == nil {
        return -1
    } else {
        if i, err := strconv.Atoi(v[0]); err != nil {
            return -1
        } else {
            return i
        }
    }
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
    v := r.GetPost(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
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

func (r *ClientRequest) GetRequestArray(k string) []string {
    v := r.GetRequest(k)
    if v == nil {
        return nil
    } else {
        return v
    }
    return nil
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

// 获取原始请求输入字符串
func (r *ClientRequest) GetJson() *gjson.Json {
    data := r.GetRaw()
    if data != nil {
        if j, err := gjson.DecodeToJson(data); err == nil {
            return j
        }
    }
    return nil
}




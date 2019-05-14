// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
    "github.com/gogf/gf/g/util/gconv"
)

// 初始化POST请求参数
func (r *Request) initPost() {
    if !r.parsedPost {
        // MultiMedia表单请求解析允许最大使用内存：1GB
        if r.ParseMultipartForm(1024*1024*1024) == nil {
            r.parsedPost = true
        }
    }
}

// 设置POST参数，仅在ghttp.Server内有效，**注意并发安全性**
func (r *Request) SetPost(key string, value string) {
    r.initPost()
    r.PostForm[key] = []string{value}
}

func (r *Request) AddPost(key string, value string) {
    r.initPost()
    r.PostForm[key] = append(r.PostForm[key], value)
}

// 获得post参数
func (r *Request) GetPost(key string, def...interface{}) []string {
    r.initPost()
    if v, ok := r.PostForm[key]; ok {
        return v
    }
    if len(def) > 0 {
        return gconv.Strings(def[0])
    }
    return nil
}

func (r *Request) GetPostString(key string, def...interface{}) string {
    value := r.GetPost(key, def...)
    if value != nil && value[0] != "" {
        return value[0]
    }
    return ""
}

func (r *Request) GetPostBool(key string, def...interface{}) bool {
    value := r.GetPostString(key, def...)
    if value != "" {
        return gconv.Bool(value)
    }
    return false
}

func (r *Request) GetPostInt(key string, def...interface{}) int {
    value := r.GetPostString(key, def...)
    if value != "" {
        return gconv.Int(value)
    }
    return 0
}

func (r *Request) GetPostInts(key string, def...interface{}) []int {
    value := r.GetPost(key, def...)
    if value != nil {
        return gconv.Ints(value)
    }
    return nil
}

func (r *Request) GetPostUint(key string, def...interface{}) uint {
    value := r.GetPostString(key, def...)
    if value != "" {
        return gconv.Uint(value)
    }
    return 0
}

func (r *Request) GetPostFloat32(key string, def...interface{}) float32 {
    value := r.GetPostString(key, def...)
    if value != "" {
        return gconv.Float32(value)
    }
    return 0
}

func (r *Request) GetPostFloat64(key string, def...interface{}) float64 {
    value := r.GetPostString(key, def...)
    if value != "" {
        return gconv.Float64(value)
    }
    return 0
}

func (r *Request) GetPostFloats(key string, def...interface{}) []float64 {
    value := r.GetPost(key, def...)
    if value != nil {
        return gconv.Floats(value)
    }
    return nil
}

func (r *Request) GetPostArray(key string, def...interface{}) []string {
    return r.GetPost(key, def...)
}

func (r *Request) GetPostStrings(key string, def...interface{}) []string {
    return r.GetPost(key, def...)
}

func (r *Request) GetPostInterfaces(key string, def...interface{}) []interface{} {
    value := r.GetPost(key, def...)
    if value != nil {
        return gconv.Interfaces(value)
    }
    return nil
}

// 获取指定键名的关联数组，并且给定当指定键名不存在时的默认值
// 需要注意的是，如果其中一个字段为数组形式，那么只会返回第一个元素，如果需要获取全部的元素，请使用GetPostArray获取特定字段内容
func (r *Request) GetPostMap(def...map[string]string) map[string]string {
    r.initPost()
    m := make(map[string]string)
    for k, v := range r.PostForm {
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

// 将所有的request参数映射到struct属性上，参数object应当为一个struct对象的指针, mapping为非必需参数，自定义参数与属性的映射关系
func (r *Request) GetPostToStruct(object interface{}, mapping...map[string]string) error {
    tagmap := r.getStructParamsTagMap(object)
    if len(mapping) > 0 {
        for k, v := range mapping[0] {
            tagmap[k] = v
        }
    }
    params := make(map[string]interface{})
    for k, v := range r.GetPostMap() {
        params[k] = v
    }
    return gconv.Struct(params, object, tagmap)
}
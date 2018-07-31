// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/util/gconv"
)

// 初始化POST请求参数
func (r *Request) initPost() {
    if !r.parsedPost.Val() {
        // 快速保存，尽量避免并发问题
        r.parsedPost.Set(true)
        // MultiMedia表单请求解析允许最大使用内存：1GB
        r.ParseMultipartForm(1024*1024*1024)
    }
}

// 设置GET参数，仅在ghttp.Server内有效，**注意并发安全性**
func (r *Request) SetPost(k string, v string) {
    r.PostForm[k] = []string{v}
}

// 获得post参数
func (r *Request) GetPost(k string) string {
    r.initPost()
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
func (r *Request) GetPostMap(defaultMap...map[string]string) map[string]string {
    r.initPost()
    m := make(map[string]string)
    if len(defaultMap) == 0 {
        for k, v := range r.PostForm {
            m[k] = v[0]
        }
    } else {
        for k, v := range defaultMap[0] {
            if v2, ok := r.PostForm[k]; ok {
                m[k] = v2[0]
            } else {
                m[k] = v
            }
        }
    }
    return m
}

// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

func (r *Request) SetRouterString(k, v string) {
    r.routerVars[k] = []string{v}
}

func (r *Request) AddRouterString(k, v string) {
    r.routerVars[k] = append(r.routerVars[k], v)
}

// 获得路由解析参数
func (r *Request) GetRouterString(k string) string {
    if v := r.GetRouterArray(k); v != nil {
        return v[0]
    }
    return ""
}

// 获得路由解析参数
func (r *Request) GetRouterArray(k string) []string {
    if v, ok := r.routerVars[k]; ok {
        return v
    }
    return nil
}


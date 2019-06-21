// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (r *Request) SetRouterString(key, value string) {
	r.routerVars[key] = []string{value}
}

func (r *Request) AddRouterString(key, value string) {
	r.routerVars[key] = append(r.routerVars[key], value)
}

// 获得路由解析参数
func (r *Request) GetRouterString(key string) string {
	if v := r.GetRouterArray(key); v != nil {
		return v[0]
	}
	return ""
}

// 获得路由解析参数
func (r *Request) GetRouterArray(key string) []string {
	if v, ok := r.routerVars[key]; ok {
		return v
	}
	return nil
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/container/gvar"

func (r *Request) SetRouterValue(key string, value interface{}) {
	r.routerMap[key] = value
}

// 获得路由解析参数
func (r *Request) GetRouterValue(key string, def ...interface{}) interface{} {
	if r.routerMap != nil {
		return r.routerMap[key]
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// 获得路由解析参数
func (r *Request) GetRouterVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRouterValue(key, def...))
}

func (r *Request) GetRouterString(key string, def ...interface{}) string {
	return r.GetRouterVar(key, def...).String()
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/container/gvar"

// 设置请求流程共享变量
func (r *Request) SetParam(key string, value interface{}) {
	if r.params == nil {
		r.params = make(map[string]interface{})
	}
	r.params[key] = value
}

// 获取请求流程共享变量
func (r *Request) GetParam(key string, def ...interface{}) interface{} {
	if r.params != nil {
		return r.params[key]
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// 获取请求流程共享变量
func (r *Request) GetParamVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetParam(key, def...))
}

// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import "gitee.com/johng/gf/g/container/gvar"

// 设置请求流程共享变量
func (r *Request) SetParam(key string, value interface{}) {
    if r.params == nil {
        r.params = make(map[string]interface{})
    }
    r.params[key] = value
}

// 获取请求流程共享变量
func (r *Request) GetParam(key string) gvar.VarRead {
    if r.params != nil {
        if v, ok := r.params[key]; ok {
            return gvar.New(v, true)
        }
    }
    return gvar.New(nil, true)
}


// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/frame/gins"
)

// 展示模板，可以给定模板参数，及临时的自定义模板函数
func (r *Response) Template(tpl string, params map[string]interface{}, funcmap...map[string]interface{}) error {
    fmap := make(gview.FuncMap)
    if len(funcmap) > 0 {
        fmap = funcmap[0]
    }
    // 内置函数
    fmap["get"]       = r.funcGet
    fmap["post"]      = r.funcPost
    fmap["request"]   = r.funcRequest
    if b, err := gins.View().Parse(tpl, params, fmap); err != nil {
        r.Write("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        r.Write(b)
    }
    return nil
}

// 模板内置函数: get
func (r *Response) funcGet(key string, def...string) gview.HTML {
    return gview.HTML(r.request.GetQueryString(key, def...))
}

// 模板内置函数: post
func (r *Response) funcPost(key string, def...string) gview.HTML {
    return gview.HTML(r.request.GetPostString(key, def...))
}

// 模板内置函数: request
func (r *Response) funcRequest(key string, def...string) gview.HTML {
    return gview.HTML(r.request.Get(key, def...))
}
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
func (r *Response) WriteTpl(tpl string, params map[string]interface{}, funcmap...map[string]interface{}) error {
    fmap := make(gview.FuncMap)
    if len(funcmap) > 0 {
        fmap = funcmap[0]
    }
    if b, err := r.ParseTpl(tpl, params, fmap); err != nil {
        r.Write("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        r.Write(b)
    }
    return nil
}

// 展示模板内容，可以给定模板参数，及临时的自定义模板函数
func (r *Response) WriteTplContent(content string, params map[string]interface{}, funcmap...map[string]interface{}) error {
    fmap := make(gview.FuncMap)
    if len(funcmap) > 0 {
        fmap = funcmap[0]
    }
    if b, err := r.ParseTplContent(content, params, fmap); err != nil {
        r.Write("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        r.Write(b)
    }
    return nil
}

// 解析模板文件，并返回模板内容
func (r *Response) ParseTpl(tpl string, params gview.Params, funcmap...map[string]interface{}) ([]byte, error) {
    fmap := make(gview.FuncMap)
    if len(funcmap) > 0 {
        fmap = funcmap[0]
    }
    return gins.View().Parse(tpl, r.buildInVars(params), r.buildInFuncs(fmap))
}

// 解析并返回模板内容
func (r *Response) ParseTplContent(content string, params gview.Params, funcmap...map[string]interface{}) ([]byte, error) {
    fmap := make(gview.FuncMap)
    if len(funcmap) > 0 {
        fmap = funcmap[0]
    }
    return gins.View().ParseContent(content, r.buildInVars(params), r.buildInFuncs(fmap))
}

// 内置变量
func (r *Response) buildInVars(params map[string]interface{}) map[string]interface{} {
    if params == nil {
        params = make(map[string]interface{})
    }

    c := gins.Config()
    if c.GetFilePath() != "" {
        params["Config"]  = c.GetMap("")
    } else {
        params["Config"]  = nil
    }
    params["Cookie"]  = r.request.Cookie.Map()
    params["Session"] = r.request.Session.Data()
    return params
}

// 内置函数
func (r *Response) buildInFuncs(funcmap map[string]interface{}) map[string]interface{} {
    if funcmap == nil {
        funcmap = make(map[string]interface{})
    }
    funcmap["get"]       = r.funcGet
    funcmap["post"]      = r.funcPost
    funcmap["request"]   = r.funcRequest
    return funcmap
}

// 模板内置函数: get
func (r *Response) funcGet(key string, def...string) string {
    return r.request.GetQueryString(key, def...)
}

// 模板内置函数: post
func (r *Response) funcPost(key string, def...string) string {
    return r.request.GetPostString(key, def...)
}

// 模板内置函数: request
func (r *Response) funcRequest(key string, def...string) string {
    return r.request.Get(key, def...)
}
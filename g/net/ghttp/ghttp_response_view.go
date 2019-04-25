// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
    "github.com/gogf/gf/g/os/gview"
    "github.com/gogf/gf/g/frame/gins"
)

// 展示模板，可以给定模板参数，及临时的自定义模板函数
func (r *Response) WriteTpl(tpl string, params...gview.Params) error {
    if b, err := r.ParseTpl(tpl, params...); err != nil {
        r.Write("Template Parsing Error: " + err.Error())
        return err
    } else {
        r.Write(b)
    }
    return nil
}

// 展示模板内容，可以给定模板参数，及临时的自定义模板函数
func (r *Response) WriteTplContent(content string, params...gview.Params) error {
    if b, err := r.ParseTplContent(content, params...); err != nil {
        r.Write("Template Parsing Error: " + err.Error())
        return err
    } else {
        r.Write(b)
    }
    return nil
}

// 解析模板文件，并返回模板内容
func (r *Response) ParseTpl(tpl string, params...gview.Params) (string, error) {
    return gins.View().Parse(tpl, r.buildInVars(params...))
}

// 解析并返回模板内容
func (r *Response) ParseTplContent(content string, params...gview.Params) (string, error) {
    return gins.View().ParseContent(content, r.buildInVars(params...))
}

// 内置变量/对象
func (r *Response) buildInVars(params...map[string]interface{}) map[string]interface{} {
	vars := map[string]interface{}(nil)
    if len(params) > 0 {
	    vars = params[0]
    } else {
	    vars = make(map[string]interface{})
    }
	vars["Config"]  = gins.Config().GetMap("")
	vars["Cookie"]  = r.request.Cookie.Map()
	vars["Session"] = r.request.Session.Map()
	vars["Get"]     = r.request.GetQueryMap()
	vars["Post"]    = r.request.GetPostMap()
    return vars
}
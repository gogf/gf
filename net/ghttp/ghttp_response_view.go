// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/util/gmode"
)

// 展示模板，可以给定模板参数，及临时的自定义模板函数
func (r *Response) WriteTpl(tpl string, params ...gview.Params) error {
	if b, err := r.ParseTpl(tpl, params...); err != nil {
		if !gmode.IsProduct() {
			r.Write("Template Parsing Error: " + err.Error())
		}
		return err
	} else {
		r.Write(b)
	}
	return nil
}

// 展示模板内容，可以给定模板参数，及临时的自定义模板函数
func (r *Response) WriteTplContent(content string, params ...gview.Params) error {
	if b, err := r.ParseTplContent(content, params...); err != nil {
		if !gmode.IsProduct() {
			r.Write("Template Parsing Error: " + err.Error())
		}
		return err
	} else {
		r.Write(b)
	}
	return nil
}

// 解析模板文件，并返回模板内容
func (r *Response) ParseTpl(tpl string, params ...gview.Params) (string, error) {
	if r.Server.config.View != nil {
		return r.Server.config.View.Parse(tpl, r.buildInVars(params...))
	}
	return gview.Instance().Parse(tpl, r.buildInVars(params...))
}

// 解析并返回模板内容
func (r *Response) ParseTplContent(content string, params ...gview.Params) (string, error) {
	if r.Server.config.View != nil {
		return r.Server.config.View.ParseContent(content, r.buildInVars(params...))
	}
	return gview.Instance().ParseContent(content, r.buildInVars(params...))
}

// 内置变量/对象
func (r *Response) buildInVars(params ...map[string]interface{}) map[string]interface{} {
	vars := map[string]interface{}(nil)
	if len(params) > 0 && params[0] != nil {
		vars = params[0]
	} else {
		vars = make(map[string]interface{})
	}
	// 当配置文件不存在时就不赋值该模板变量，不然会报错
	if c := gins.Config(); c.FilePath() != "" {
		vars["Config"] = c.GetMap(".")
	}
	vars["Get"] = r.Request.GetQueryMap()
	vars["Post"] = r.Request.GetPostMap()
	vars["Cookie"] = r.Request.Cookie.Map()
	vars["Session"] = r.Request.Session.Map()
	return vars
}

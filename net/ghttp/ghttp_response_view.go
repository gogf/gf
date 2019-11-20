// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/util/gmode"
)

// SetView sets template view engine object for this response.
func (r *Response) SetView(view *gview.View) {
	r.view = view
}

// GetView returns the template view engine object for this response.
func (r *Response) GetView() *gview.View {
	view := r.view
	if view == nil {
		view = r.Server.config.View
	}
	if view == nil {
		gview.Instance()
	}
	return view
}

// WriteTpl parses and responses given template file.
// The parameter <params> specifies the template variables for parsing.
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

// WriteTplDefault parses and responses the default template file.
// The parameter <params> specifies the template variables for parsing.
func (r *Response) WriteTplDefault(params ...gview.Params) error {
	if b, err := r.ParseTplDefault(params...); err != nil {
		if !gmode.IsProduct() {
			r.Write("Template Parsing Error: " + err.Error())
		}
		return err
	} else {
		r.Write(b)
	}
	return nil
}

// WriteTplContent parses and responses the template content.
// The parameter <params> specifies the template variables for parsing.
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

// ParseTpl parses given template file <tpl> with given template variables <params>
// and returns the parsed template content.
func (r *Response) ParseTpl(tpl string, params ...gview.Params) (string, error) {
	return r.GetView().Parse(tpl, r.buildInVars(params...))
}

// ParseDefault parses the default template file with params.
func (r *Response) ParseTplDefault(params ...gview.Params) (string, error) {
	return r.GetView().ParseDefault(r.buildInVars(params...))
}

// ParseTplContent parses given template file <file> with given template parameters <params>
// and returns the parsed template content.
func (r *Response) ParseTplContent(content string, params ...gview.Params) (string, error) {
	return r.GetView().ParseContent(content, r.buildInVars(params...))
}

// buildInVars merges build-in variables into <params> and returns the new template variables.
func (r *Response) buildInVars(params ...map[string]interface{}) map[string]interface{} {
	var vars map[string]interface{}
	if len(params) > 0 && params[0] != nil {
		vars = params[0]
	} else {
		vars = make(map[string]interface{})
	}
	// 当配置文件不存在时就不赋值该模板变量，不然会报错
	if c := gcfg.Instance(); c.FilePath() != "" {
		vars["Config"] = c.GetMap(".")
	}
	vars["Get"] = r.Request.GetQueryMap()
	vars["Post"] = r.Request.GetPostMap()
	vars["Cookie"] = r.Request.Cookie.Map()
	vars["Session"] = r.Request.Session.Map()
	return vars
}

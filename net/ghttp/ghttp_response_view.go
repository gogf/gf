// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package ghttp

import (
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmode"
	"github.com/gogf/gf/util/gutil"
)

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
	return r.Request.GetView().Parse(r.Request.Context(), tpl, r.buildInVars(params...))
}

// ParseTplDefault parses the default template file with params.
func (r *Response) ParseTplDefault(params ...gview.Params) (string, error) {
	return r.Request.GetView().ParseDefault(r.Request.Context(), r.buildInVars(params...))
}

// ParseTplContent parses given template file <file> with given template parameters <params>
// and returns the parsed template content.
func (r *Response) ParseTplContent(content string, params ...gview.Params) (string, error) {
	return r.Request.GetView().ParseContent(r.Request.Context(), content, r.buildInVars(params...))
}

// buildInVars merges build-in variables into <params> and returns the new template variables.
// TODO performance improving.
func (r *Response) buildInVars(params ...map[string]interface{}) map[string]interface{} {
	m := gutil.MapMergeCopy(r.Request.viewParams)
	if len(params) > 0 {
		gutil.MapMerge(m, params[0])
	}
	// Retrieve custom template variables from request object.
	sessionMap := gconv.MapDeep(r.Request.Session.Map())
	gutil.MapMerge(m, map[string]interface{}{
		"Form":    r.Request.GetFormMap(),
		"Query":   r.Request.GetQueryMap(),
		"Request": r.Request.GetMap(),
		"Cookie":  r.Request.Cookie.Map(),
		"Session": sessionMap,
	})
	// Note that it should assign no Config variable to template
	// if there's no configuration file.
	if c := gcfg.Instance(); c.Available() {
		m["Config"] = c.Map()
	}
	return m
}

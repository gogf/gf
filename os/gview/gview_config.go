// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

import "github.com/gogf/gf/i18n/gi18n"

// Assigns binds multiple global template variables to current view object.
// Note that it's not concurrent-safe, which means it would panic
// if it's called in multiple goroutines in runtime.
func (view *View) Assigns(data Params) {
	for k, v := range data {
		view.data[k] = v
	}
}

// Assign binds a global template variable to current view object.
// Note that it's not concurrent-safe, which means it would panic
// if it's called in multiple goroutines in runtime.
func (view *View) Assign(key string, value interface{}) {
	view.data[key] = value
}

// SetDefaultFile sets default template file for parsing.
func (view *View) SetDefaultFile(file string) {
	view.defaultFile = file
}

// SetDelimiters sets customized delimiters for template parsing.
func (view *View) SetDelimiters(left, right string) {
	view.delimiters[0] = left
	view.delimiters[1] = right
}

// BindFunc registers customized global template function named <name>
// with given function <function> to current view object.
// The <name> is the function name which can be called in template content.
func (view *View) BindFunc(name string, function interface{}) {
	view.funcMap[name] = function
}

// BindFuncMap registers customized global template functions by map to current view object.
// The key of map is the template function name
// and the value of map is the address of customized function.
func (view *View) BindFuncMap(funcMap FuncMap) {
	for k, v := range funcMap {
		view.funcMap[k] = v
	}
}

// SetI18n binds i18n manager to current view engine.
func (view *View) SetI18n(manager *gi18n.Manager) {
	view.i18nManager = manager
}

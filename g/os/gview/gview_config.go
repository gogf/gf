// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview

// Assign binds multiple template variables to current view object.
// Each goroutine will take effect after the call, so it is concurrent-safe.
func (view *View) Assigns(data Params) {
	view.mu.Lock()
	for k, v := range data {
		view.data[k] = v
	}
	view.mu.Unlock()
}

// Assign binds a template variable to current view object.
// Each goroutine will take effect after the call, so it is concurrent-safe.
func (view *View) Assign(key string, value interface{}) {
	view.mu.Lock()
	view.data[key] = value
	view.mu.Unlock()
}

// SetDelimiters sets customized delimiters for template parsing.
func (view *View) SetDelimiters(left, right string) {
	view.mu.Lock()
    view.delimiters[0] = left
    view.delimiters[1] = right
	view.mu.Unlock()
}

// BindFunc registers customized template function named <name>
// with given function <function> to current view object.
// The <name> is the function name which can be called in template content.
func (view *View) BindFunc(name string, function interface{}) {
    view.mu.Lock()
    view.funcMap[name] = function
    view.mu.Unlock()
}

// BindFuncMap registers customized template functions by map to current view object.
// The key of map is the template function name
// and the value of map is the address of customized function.
func (view *View) BindFuncMap(funcMap FuncMap) {
	view.mu.Lock()
	for k, v := range funcMap {
		view.funcMap[k] = v
	}
	view.mu.Unlock()
}

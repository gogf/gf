// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmvc

import (
	"github.com/gogf/gf/frame/gins"
	"sync"

	"github.com/gogf/gf/util/gmode"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gview"
)

// View is the view object for controller.
// It's initialized when controller request initializes and destroyed
// when the controller request closes.
// Deprecated, no longer suggested.
type View struct {
	mu       sync.RWMutex
	view     *gview.View
	data     gview.Params
	response *ghttp.Response
}

// NewView creates and returns a controller view object.
// Deprecated, no longer suggested.
func NewView(w *ghttp.Response) *View {
	return &View{
		view:     gins.View(),
		data:     make(gview.Params),
		response: w,
	}
}

// Assigns assigns template variables to this view object.
func (view *View) Assigns(data gview.Params) {
	view.mu.Lock()
	for k, v := range data {
		view.data[k] = v
	}
	view.mu.Unlock()
}

// Assign assigns one template variable to this view object.
func (view *View) Assign(key string, value interface{}) {
	view.mu.Lock()
	view.data[key] = value
	view.mu.Unlock()
}

// Parse parses given template file <tpl> with assigned template variables
// and returns the parsed template content.
func (view *View) Parse(file string) (string, error) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	buffer, err := view.response.ParseTpl(file, view.data)
	return buffer, err
}

// ParseContent parses given template file <file> with assigned template variables
// and returns the parsed template content.
func (view *View) ParseContent(content string) (string, error) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	buffer, err := view.response.ParseTplContent(content, view.data)
	return buffer, err
}

// LockFunc locks writing for template variables by callback function <f>.
func (view *View) LockFunc(f func(data gview.Params)) {
	view.mu.Lock()
	defer view.mu.Unlock()
	f(view.data)
}

// RLockFunc locks reading for template variables by callback function <f>.
func (view *View) RLockFunc(f func(data gview.Params)) {
	view.mu.RLock()
	defer view.mu.RUnlock()
	f(view.data)
}

// BindFunc registers customized template function named <name>
// with given function <function> to current view object.
// The <name> is the function name which can be called in template content.
func (view *View) BindFunc(name string, function interface{}) {
	view.view.BindFunc(name, function)
}

// BindFuncMap registers customized template functions by map to current view object.
// The key of map is the template function name
// and the value of map is the address of customized function.
func (view *View) BindFuncMap(funcMap gview.FuncMap) {
	view.view.BindFuncMap(funcMap)
}

// Display parses and writes the parsed template file content to http response.
func (view *View) Display(file ...string) error {
	name := view.view.GetDefaultFile()
	if len(file) > 0 {
		name = file[0]
	}
	if content, err := view.Parse(name); err != nil {
		if !gmode.IsProduct() {
			view.response.Write("Tpl Parsing Error: " + err.Error())
		}
		return err
	} else {
		view.response.Write(content)
	}
	return nil
}

// DisplayContent parses and writes the parsed content to http response.
func (view *View) DisplayContent(content string) error {
	if content, err := view.ParseContent(content); err != nil {
		if !gmode.IsProduct() {
			view.response.Write("Tpl Parsing Error: " + err.Error())
		}
		return err
	} else {
		view.response.Write(content)
	}
	return nil
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
)

// AdapterContent implements interface Adapter using content.
// The configuration content supports the coding types as package `gjson`.
type AdapterContent struct {
	jsonVar *gvar.Var // The pared JSON object for configuration content, type: *gjson.Json.
}

// NewAdapterContent returns a new configuration management object using custom content.
// The parameter `content` specifies the default configuration content for reading.
func NewAdapterContent(content ...string) (*AdapterContent, error) {
	a := &AdapterContent{
		jsonVar: gvar.New(nil, true),
	}
	if len(content) > 0 {
		if err := a.SetContent(content[0]); err != nil {
			return nil, err
		}
	}
	return a, nil
}

// SetContent sets customized configuration content for specified `file`.
// The `file` is unnecessary param, default is DefaultConfigFile.
func (a *AdapterContent) SetContent(content string) error {
	j, err := gjson.LoadContent(content, true)
	if err != nil {
		return gerror.Wrap(err, `load configuration content failed`)
	}
	a.jsonVar.Set(j)
	return nil
}

// Available checks and returns the backend configuration service is available.
// The optional parameter `resource` specifies certain configuration resource.
//
// Note that this function does not return error as it just does simply check for
// backend configuration service.
func (a *AdapterContent) Available(ctx context.Context, resource ...string) (ok bool) {
	if a.jsonVar.IsNil() {
		return false
	}
	return true
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (a *AdapterContent) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Val().(*gjson.Json).Get(pattern).Val(), nil
}

// Data retrieves and returns all configuration data in current resource as map.
// Note that this function may lead lots of memory usage if configuration data is too large,
// you can implement this function if necessary.
func (a *AdapterContent) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if a.jsonVar.IsNil() {
		return nil, nil
	}
	return a.jsonVar.Val().(*gjson.Json).Var().Map(), nil
}

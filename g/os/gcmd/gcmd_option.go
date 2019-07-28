// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import "github.com/gogf/gf/g/container/gvar"

// GetAll returns all option values as map[string]string.
func (c *gCmdOption) GetAll() map[string]string {
	return c.options
}

// BuildOptions builds the options as string.
func (c *gCmdOption) Build(prefix ...string) string {
	return BuildOptions(c.options, prefix...)
}

// Get returns the option value string specified by <key>,
// if value dose not exist, then returns <def>.
func (c *gCmdOption) Get(key string, def ...string) string {
	if option, ok := c.options[key]; ok {
		return option
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

// Contains checks whether the option named <key> exists.
func (c *gCmdOption) Contains(key string) bool {
	_, ok := c.options[key]
	return ok
}

// Set sets the option named <key> with value <value>.
func (c *gCmdOption) Set(key string, value string) {
	c.options[key] = value
}

// GetVar returns the option value as *gvar.Var object specified by <key>,
// if value does not exist, then returns <def> as its default value.
func (c *gCmdOption) GetVar(key string, def ...string) *gvar.Var {
	return gvar.New(c.Get(key, def...))
}

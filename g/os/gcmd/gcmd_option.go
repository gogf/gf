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

// Get returns the option value string specified by <key>,
// if value dose not exist, then returns <def>.
func (c *gCmdOption) Get(key string, def...string) string {
    if option, ok := c.options[key]; ok {
        return option
    } else if len(def) > 0 {
        return def[0]
    }
    return ""
}

// GetVar returns the option value as *gvar.Var object specified by <key>,
// if value does not exist, then returns <def> as its default value.
func (c *gCmdOption) GetVar(key string, def...string) *gvar.Var {
	return gvar.New(c.Get(key, def...), true)
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import "github.com/gogf/gf/g/container/gvar"

// GetAll returns all values as a slice of string.
func (c *gCmdValue) GetAll() []string {
    return c.values
}

// Get returns value by index <index> as string,
// if value does not exist, then returns <def>.
func (c *gCmdValue) Get(index int, def...string) string {
    if index < len(c.values) {
        return c.values[index]
    } else if len(def) > 0 {
        return def[0]
    }
    return ""
}

// GetVar returns value by index <index> as *gvar.Var object,
// if value does not exist, then returns <def> as its default value.
func (c *gCmdValue) GetVar(index int, def...string) *gvar.Var {
	return gvar.New(c.Get(index, def...), true)
}

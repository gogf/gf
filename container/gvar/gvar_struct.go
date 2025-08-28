// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Struct maps value of `v` to `pointer`.
// The parameter `pointer` should be a pointer to a struct instance.
// The parameter `mapping` is used to specify the key-to-attribute mapping rules.
func (v *Var) Struct(pointer any, mapping ...map[string]string) error {
	return gconv.Struct(v.Val(), pointer, mapping...)
}

// Structs converts and returns `v` as given struct slice.
func (v *Var) Structs(pointer any, mapping ...map[string]string) error {
	return gconv.Structs(v.Val(), pointer, mapping...)
}

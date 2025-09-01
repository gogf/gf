// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Scan automatically checks the type of `pointer` and converts value of Var to `pointer`.
//
// See gconv.Scan.
func (v *Var) Scan(pointer any, mapping ...map[string]string) error {
	return gconv.Scan(v.Val(), pointer, mapping...)
}

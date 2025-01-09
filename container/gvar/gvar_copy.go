// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/internal/deepcopy"
	"github.com/gogf/gf/v2/util/gutil"
)

// Copy does a deep copy of current Var and returns a pointer to this Var.
func (v *Var) Copy() *Var {
	return New(gutil.Copy(v.Val()), v.safe)
}

// Clone does a shallow copy of current Var and returns a pointer to this Var.
func (v *Var) Clone() *Var {
	return New(v.Val(), v.safe)
}

// DeepCopy implements interface for deep copy of current type.
func (v *Var) DeepCopy() interface{} {
	if v == nil {
		return nil
	}
	return New(deepcopy.Copy(v.Val()), v.safe)
}

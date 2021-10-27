// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/internal/utils"
)

// IsNil checks whether `v` is nil.
func (v *Var) IsNil() bool {
	return utils.IsNil(v.Val())
}

// IsEmpty checks whether `v` is empty.
func (v *Var) IsEmpty() bool {
	return utils.IsEmpty(v.Val())
}

// IsInt checks whether `v` is type of int.
func (v *Var) IsInt() bool {
	return utils.IsInt(v.Val())
}

// IsUint checks whether `v` is type of uint.
func (v *Var) IsUint() bool {
	return utils.IsUint(v.Val())
}

// IsFloat checks whether `v` is type of float.
func (v *Var) IsFloat() bool {
	return utils.IsFloat(v.Val())
}

// IsSlice checks whether `v` is type of slice.
func (v *Var) IsSlice() bool {
	return utils.IsSlice(v.Val())
}

// IsMap checks whether `v` is type of map.
func (v *Var) IsMap() bool {
	return utils.IsMap(v.Val())
}

// IsStruct checks whether `v` is type of struct.
func (v *Var) IsStruct() bool {
	return utils.IsStruct(v.Val())
}

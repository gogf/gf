// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"github.com/gogf/gf/v2/internal/deepcopy"
)

// Copy returns a deep copy of v.
//
// Copy is unable to copy unexported fields in a struct (lowercase field names).
// Unexported fields can't be reflected by the Go runtime and therefore
// they can't perform any data copies.
func Copy(src interface{}) (dst interface{}) {
	return deepcopy.Copy(src)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

// Any is a struct for concurrent-safe operation for type any.
type Any = Interface

// NewAny creates and returns a concurrent-safe object for any type,
// with given initial value `value`.
func NewAny(value ...any) *Any {
	t := &Any{}
	if len(value) > 0 && value[0] != nil {
		t.value.Store(value[0])
	}
	return t
}

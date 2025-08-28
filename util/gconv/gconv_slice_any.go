// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceAny is alias of Interfaces.
func SliceAny(any any) []any {
	return Interfaces(any)
}

// Interfaces converts `any` to []any.
func Interfaces(any any) []any {
	result, _ := defaultConverter.SliceAny(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}

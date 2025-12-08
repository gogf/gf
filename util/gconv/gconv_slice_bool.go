// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceBool is alias of Bools.
func SliceBool(anyInput any) []bool {
	return Bools(anyInput)
}

// Bools converts `any` to []bool.
func Bools(anyInput any) []bool {
	result, _ := defaultConverter.SliceBool(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}

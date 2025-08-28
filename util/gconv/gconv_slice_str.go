// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceStr is alias of Strings.
func SliceStr(anyInput any) []string {
	return Strings(anyInput)
}

// Strings converts `any` to []string.
func Strings(anyInput any) []string {
	result, _ := defaultConverter.SliceStr(anyInput, SliceOption{
		ContinueOnError: true,
	})
	return result
}

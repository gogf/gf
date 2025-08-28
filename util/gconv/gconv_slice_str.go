// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// SliceStr is alias of Strings.
func SliceStr(any any) []string {
	return Strings(any)
}

// Strings converts `any` to []string.
func Strings(any any) []string {
	result, _ := defaultConverter.SliceStr(any, SliceOption{
		ContinueOnError: true,
	})
	return result
}

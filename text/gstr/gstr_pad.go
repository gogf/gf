// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

// PadLeft this string with another string (multiple times, if needed)
// until the resulting string reaches the given length. The padding
// is applied from the start of this string.
func PadLeft(s, padStr string, targetLen int) string {
	if len(s) < targetLen {
		for i := len(s); i < targetLen; i++ {
			s = padStr + s
		}
	}
	return s
}

// PadRight this string with another string (multiple times, if needed)
// until the resulting string reaches the given length. The padding
// is applied from the end of this string.
func PadRight(s, padStr string, targetLen int) string {
	if len(s) < targetLen {
		for i := len(s); i < targetLen; i++ {
			s += padStr
		}
	}
	return s
}

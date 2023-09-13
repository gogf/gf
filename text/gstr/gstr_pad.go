// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

// pads this string with another string (multiple times, if needed)
// until the resulting string reaches the given length. The padding
// is applied from the start of this string.
func PadStart(s, padStr string, targetLen int) string {
	if len(s) < targetLen {
		for i := len(s); i < targetLen; i++ {
			s = padStr + s
		}
	}
	return s
}

// pads this string with another string (multiple times, if needed)
// until the resulting string reaches the given length. The padding
// is applied from the end of this string.
func PadEnd(s, padStr string, targetLen int) string {
	if len(s) < targetLen {
		for i := len(s); i < targetLen; i++ {
			s += padStr
		}
	}
	return s
}

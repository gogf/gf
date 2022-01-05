// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Repeat returns a new string consisting of multiplier copies of the string input.
func Repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

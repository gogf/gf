// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Compare returns an integer comparing two strings lexicographically.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

// Equal reports whether `a` and `b`, interpreted as UTF-8 strings,
// are equal under Unicode case-folding, case-insensitively.
func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

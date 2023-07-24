// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "unicode/utf8"

// LenRune returns string length of unicode.
func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}

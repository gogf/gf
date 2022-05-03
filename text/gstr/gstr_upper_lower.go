// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"

	"strings"

	"github.com/gogf/gf/v2/internal/utils"
)

// ToLower returns a copy of the string s with all Unicode letters mapped to their lower case.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper returns a copy of the string s with all Unicode letters mapped to their upper case.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// UcFirst returns a copy of the string s with the first letter mapped to its upper case.
func UcFirst(s string) string {
	return utils.UcFirst(s)
}

// LcFirst returns a copy of the string s with the first letter mapped to its lower case.
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if IsLetterUpper(s[0]) {
		return string(s[0]+32) + s[1:]
	}
	return s
}

// UcWords uppercase the first character of each word in a string.
func UcWords(str string) string {
	return strings.Title(str)
}

// IsLetterLower tests whether the given byte b is in lower case.
func IsLetterLower(b byte) bool {
	return utils.IsLetterLower(b)
}

// IsLetterUpper tests whether the given byte b is in upper case.
func IsLetterUpper(b byte) bool {
	return utils.IsLetterUpper(b)
}

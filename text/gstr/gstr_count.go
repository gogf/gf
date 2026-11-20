// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"
	"unicode"
)

// Count counts the number of `substr` appears in `s`.
// It returns 0 if no `substr` found in `s`.
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// CountI counts the number of `substr` appears in `s`, case-insensitively.
// It returns 0 if no `substr` found in `s`.
func CountI(s, substr string) int {
	return strings.Count(ToLower(s), ToLower(substr))
}

// CountWords returns information about words' count used in a string.
// It considers parameter `str` as unicode string.
func CountWords(str string) map[string]int {
	m := make(map[string]int)
	for _, word := range strings.Fields(str) {
		m[word]++
	}
	return m
}

// CountChars returns information about chars' count used in a string.
// It considers parameter `str` as Unicode string.
func CountChars(str string, noSpace ...bool) map[string]int {
	m := make(map[string]int)
	countSpace := len(noSpace) == 0 || !noSpace[0]
	for _, r := range str {
		if !countSpace && unicode.IsSpace(r) {
			continue
		}
		m[string(r)]++
	}
	return m
}

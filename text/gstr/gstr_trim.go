// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"github.com/gogf/gf/internal/utils"
	"strings"
)

// Trim strips whitespace (or other characters) from the beginning and end of a string.
// The optional parameter <characterMask> specifies the additional stripped characters.
func Trim(str string, characterMask ...string) string {
	return utils.Trim(str, characterMask...)
}

// TrimStr strips all the given <cut> string from the beginning and end of a string.
// Note that it does not strip the whitespaces of its beginning or end.
func TrimStr(str string, cut string, count ...int) string {
	return TrimLeftStr(TrimRightStr(str, cut, count...), cut, count...)
}

// TrimLeft strips whitespace (or other characters) from the beginning of a string.
func TrimLeft(str string, characterMask ...string) string {
	trimChars := utils.DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.TrimLeft(str, trimChars)
}

// TrimLeftStr strips all the given <cut> string from the beginning of a string.
// Note that it does not strip the whitespaces of its beginning.
func TrimLeftStr(str string, cut string, count ...int) string {
	var (
		lenCut   = len(cut)
		cutCount = 0
	)
	for len(str) >= lenCut && str[0:lenCut] == cut {
		str = str[lenCut:]
		cutCount++
		if len(count) > 0 && count[0] != -1 && cutCount >= count[0] {
			break
		}
	}
	return str
}

// TrimRight strips whitespace (or other characters) from the end of a string.
func TrimRight(str string, characterMask ...string) string {
	trimChars := utils.DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.TrimRight(str, trimChars)
}

// TrimRightStr strips all the given <cut> string from the end of a string.
// Note that it does not strip the whitespaces of its end.
func TrimRightStr(str string, cut string, count ...int) string {
	var (
		lenStr   = len(str)
		lenCut   = len(cut)
		cutCount = 0
	)
	for lenStr >= lenCut && str[lenStr-lenCut:lenStr] == cut {
		lenStr = lenStr - lenCut
		str = str[:lenStr]
		cutCount++
		if len(count) > 0 && count[0] != -1 && cutCount >= count[0] {
			break
		}
	}
	return str
}

// TrimAll trims all characters in string `str`.
func TrimAll(str string, characterMask ...string) string {
	trimChars := utils.DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	var (
		filtered bool
		slice    = make([]rune, 0, len(str))
	)
	for _, char := range str {
		filtered = false
		for _, trimChar := range trimChars {
			if char == trimChar {
				filtered = true
				break
			}
		}
		if !filtered {
			slice = append(slice, char)
		}
	}
	return string(slice)
}

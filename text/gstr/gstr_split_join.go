// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gconv"
)

// Split splits string `str` by a string `delimiter`, to an array.
func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

// SplitAndTrim splits string `str` by a string `delimiter` to an array,
// and calls Trim to every element of this array. It ignores the elements
// which are empty after Trim.
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	return utils.SplitAndTrim(str, delimiter, characterMask...)
}

// Join concatenates the elements of `array` to create a single string. The separator string
// `sep` is placed between elements in the resulting string.
func Join(array []string, sep string) string {
	return strings.Join(array, sep)
}

// JoinAny concatenates the elements of `array` to create a single string. The separator string
// `sep` is placed between elements in the resulting string.
//
// The parameter `array` can be any type of slice, which be converted to string array.
func JoinAny(array interface{}, sep string) string {
	return strings.Join(gconv.Strings(array), sep)
}

// Explode splits string `str` by a string `delimiter`, to an array.
// See http://php.net/manual/en/function.explode.php.
func Explode(delimiter, str string) []string {
	return Split(str, delimiter)
}

// Implode joins array elements `pieces` with a string `glue`.
// http://php.net/manual/en/function.implode.php
func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

// ChunkSplit splits a string into smaller chunks.
// Can be used to split a string into smaller chunks which is useful for
// e.g. converting BASE64 string output to match RFC 2045 semantics.
// It inserts end every chunkLen characters.
// It considers parameter `body` and `end` as unicode string.
func ChunkSplit(body string, chunkLen int, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, endRunes := []rune(body), []rune(end)
	l := len(runes)
	if l <= 1 || l < chunkLen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(endRunes))
	for i := 0; i < l; i += chunkLen {
		if i+chunkLen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunkLen]...)
		}
		ns = append(ns, endRunes...)
	}
	return string(ns)
}

// Fields returns the words used in a string as slice.
func Fields(str string) []string {
	return strings.Fields(str)
}

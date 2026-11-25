// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"

	"github.com/gogf/gf/v2/internal/utils"
)

// Replace returns a copy of the string `origin`
// in which string `search` replaced by `replace` case-sensitively.
func Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}

// ReplaceI returns a copy of the string `origin`
// in which string `search` replaced by `replace` case-insensitively.
// It returns `origin` unchanged if `search` is empty.
func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 || search == "" {
		return origin
	}
	var result strings.Builder
	for ; n != 0; n-- {
		// The match is cut by its byte range in origin, which can differ in length from search.
		start, end := indexFold(origin, search, false)
		if start == -1 {
			break
		}
		result.WriteString(origin[:start])
		result.WriteString(replace)
		origin = origin[end:]
	}
	result.WriteString(origin)
	return result.String()
}

// ReplaceByArray returns a copy of `origin`,
// which is replaced by a slice in order, case-sensitively.
func ReplaceByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = Replace(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceIByArray returns a copy of `origin`,
// which is replaced by a slice in order, case-insensitively.
func ReplaceIByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = ReplaceI(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceByMap returns a copy of `origin`,
// which is replaced by a map in unordered way, case-sensitively.
func ReplaceByMap(origin string, replaces map[string]string) string {
	return utils.ReplaceByMap(origin, replaces)
}

// ReplaceIByMap returns a copy of `origin`,
// which is replaced by a map in unordered way, case-insensitively.
func ReplaceIByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = ReplaceI(origin, k, v)
	}
	return origin
}

// ReplaceFunc returns a copy of the string `origin` in which each non-overlapping substring
// that matches the given search string is replaced by the result of function `f` applied to that substring.
// The function `f` is called with each matching substring as its argument and must return a string to be used
// as the replacement value.
func ReplaceFunc(origin string, search string, f func(string) string) string {
	if search == "" {
		return origin
	}
	var (
		searchLen = len(search)
		originLen = len(origin)
	)
	// If search string is longer than origin string, no match is possible
	if searchLen > originLen {
		return origin
	}
	var (
		result     strings.Builder
		lastMatch  int
		currentPos int
	)
	// Pre-allocate the builder capacity to avoid reallocations
	result.Grow(originLen)

	for currentPos < originLen {
		pos := Pos(origin[currentPos:], search)
		if pos == -1 {
			break
		}
		pos += currentPos
		// Append unmatched portion
		result.WriteString(origin[lastMatch:pos])
		// Apply replacement function and append result
		match := origin[pos : pos+searchLen]
		result.WriteString(f(match))
		// Update positions
		lastMatch = pos + searchLen
		currentPos = lastMatch
	}
	// Append remaining unmatched portion
	if lastMatch < originLen {
		result.WriteString(origin[lastMatch:])
	}
	return result.String()
}

// ReplaceIFunc returns a copy of the string `origin` in which each non-overlapping substring
// that matches the given search string is replaced by the result of function `f` applied to that substring.
// The match is done case-insensitively.
// The function `f` is called with each matching substring as its argument and must return a string to be used
// as the replacement value.
func ReplaceIFunc(origin string, search string, f func(string) string) string {
	if search == "" {
		return origin
	}
	var (
		originLen  = len(origin)
		result     strings.Builder
		lastMatch  int
		currentPos int
	)
	// Pre-allocate the builder capacity to avoid reallocations
	result.Grow(originLen)

	for currentPos < originLen {
		// The match is cut by its byte range in origin, which can differ in length from search.
		pos, end := indexFold(origin[currentPos:], search, false)
		if pos == -1 {
			break
		}
		pos += currentPos
		end += currentPos
		// Append unmatched portion
		result.WriteString(origin[lastMatch:pos])
		// Apply replacement function and append result
		result.WriteString(f(origin[pos:end]))
		// Update positions
		lastMatch = end
		currentPos = lastMatch
	}
	// Append remaining unmatched portion
	if lastMatch < originLen {
		result.WriteString(origin[lastMatch:])
	}
	return result.String()
}

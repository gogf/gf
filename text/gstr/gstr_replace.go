// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"github.com/gogf/gf/v2/internal/utils"
	"strings"
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
func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return origin
	}
	var (
		length      = len(search)
		searchLower = strings.ToLower(search)
	)
	for {
		originLower := strings.ToLower(origin)
		if pos := strings.Index(originLower, searchLower); pos != -1 {
			origin = origin[:pos] + replace + origin[pos+length:]
			if n--; n == 0 {
				break
			}
		} else {
			break
		}
	}
	return origin
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

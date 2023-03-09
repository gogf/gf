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

// IsNumeric tests whether the given string s is numeric.
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}

// StartsWith Determine if `haystack` starts with `needle`.
func StartsWith(haystack, needle string) bool {
	length := len(needle)
	if len(haystack) < length {
		return false
	}
	return haystack[:length] == needle
}

// StartsWithI Determine if `haystack` starts with `needle`.
func StartsWithI(haystack, needle string) bool {
	length := len(needle)
	if len(haystack) < length {
		return false
	}

	return strings.ToLower(haystack[:length]) == strings.ToLower(needle)
}

// EndsWith Determine if `haystack` ends with `needle`.
func EndsWith(haystack, needle string) bool {
	length := len(needle)
	if len(haystack) < length {
		return false
	}
	return haystack[len(haystack)-length:] == needle
}

// EndsWithI Determine if `haystack` ends with `needle`.
func EndsWithI(haystack, needle string) bool {
	length := len(needle)
	if len(haystack) < length {
		return false
	}
	return strings.ToLower(haystack[len(haystack)-length:]) == strings.ToLower(needle)
}

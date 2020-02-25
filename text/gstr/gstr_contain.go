// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Contains reports whether <substr> is within <str>, case-sensitively.
func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ContainsI reports whether substr is within str, case-insensitively.
func ContainsI(str, substr string) bool {
	return PosI(str, substr) != -1
}

// ContainsAny reports whether any Unicode code points in <chars> are within <s>.
func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}

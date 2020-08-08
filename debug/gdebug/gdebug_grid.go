// Copyright 2019-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"regexp"
	"runtime"
	"strconv"
)

var (
	// gridRegex is the regular expression object for parsing goroutine id from stack information.
	gridRegex = regexp.MustCompile(`^\w+\s+(\d+)\s+`)
)

// GoroutineId retrieves and returns the current goroutine id from stack information.
// Be very aware that, it is with low performance as it uses runtime.Stack function.
// It is commonly used for debugging purpose.
func GoroutineId() int {
	buf := make([]byte, 26)
	runtime.Stack(buf, false)
	match := gridRegex.FindSubmatch(buf)
	id, _ := strconv.Atoi(string(match[1]))
	return id
}

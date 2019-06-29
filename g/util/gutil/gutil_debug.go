// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"github.com/gogf/gf/g/internal/debug"
)

// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func PrintStack(skip ...int) {
	number := 1
	if len(skip) > 0 {
		number = skip[0] + 1
	}
	debug.PrintStack(number)
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack(skip ...int) string {
	number := 1
	if len(skip) > 0 {
		number = skip[0] + 1
	}
	return debug.Stack(number)
}

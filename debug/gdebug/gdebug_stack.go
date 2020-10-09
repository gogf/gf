// Copyright 2019-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func PrintStack(skip ...int) {
	fmt.Print(Stack(skip...))
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack(skip ...int) string {
	return StackWithFilter("", skip...)
}

// StackWithFilter returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter <filter> is used to filter the path of the caller.
func StackWithFilter(filter string, skip ...int) string {
	return StackWithFilters([]string{filter}, skip...)
}

// StackWithFilters returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter <filters> is a slice of strings, which are used to filter the path of the
// caller.
//
// TODO Improve the performance using debug.Stack.
func StackWithFilters(filters []string, skip ...int) string {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	var (
		name                  = ""
		space                 = "  "
		index                 = 1
		buffer                = bytes.NewBuffer(nil)
		filtered              = false
		ok                    = true
		pc, file, line, start = callerFromIndex(filters)
	)
	for i := start + number; i < gMAX_DEPTH; i++ {
		if i != start {
			pc, file, line, ok = runtime.Caller(i)
		}
		if ok {
			// Filter empty file.
			if file == "" {
				continue
			}
			// GOROOT filter.
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			// Custom filtering.
			filtered = false
			for _, filter := range filters {
				if filter != "" && strings.Contains(file, filter) {
					filtered = true
					break
				}
			}
			if filtered {
				continue
			}
			if strings.Contains(file, gFILTER_KEY) {
				continue
			}
			if fn := runtime.FuncForPC(pc); fn == nil {
				name = "unknown"
			} else {
				name = fn.Name()
			}
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf("%d.%s%s\n    %s:%d\n", index, space, name, file, line))
			index++
		} else {
			break
		}
	}
	return buffer.String()
}

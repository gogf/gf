// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package debug contains facilities for programs to debug themselves while
// they are running.
package debug

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func PrintStack(skip ...int) {
	fmt.Print(string(Stack(skip...)))
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack(skip ...int) []byte {
	buffer := make([]byte, 512)
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	for {
		n := runtime.Stack(buffer, false)
		if n < len(buffer) {
			lines := bytes.Split(buffer[:n], []byte{'\n'})
			index := 1
			stacks := bytes.NewBuffer(nil)
			for i, line := range lines {
				if i < 5+number*2 || len(line) == 0 {
					continue
				}
				if i%2 != 0 {
					stacks.WriteString(strconv.Itoa(index) + ".\t")
					index++
				}
				stacks.Write(line)
				stacks.WriteByte('\n')
			}
			return stacks.Bytes()
		}
		buffer = make([]byte, 2*len(buffer))
	}
}

// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import "runtime"

// stack represents a stack of program counters.
type stack []uintptr

const (
	gMAX_STACK_DEPTH = 32
)

func callers() stack {
	var pcs [gMAX_STACK_DEPTH]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}

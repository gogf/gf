// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Func is the function which contains context parameter.
type Func func(ctx context.Context)

// RecoverFunc is the panic recover function which contains context parameter.
type RecoverFunc func(ctx context.Context, exception error)

// Go creates a new asynchronous goroutine function with specified recover function.
//
// The parameter `recoverFunc` is called when any panic during executing of `goroutineFunc`.
// If `recoverFunc` is not given or given nil, it ignores the panic from `goroutineFunc`.
func Go(ctx context.Context, goroutineFunc Func, recoverFunc RecoverFunc) {
	if goroutineFunc == nil {
		return
	}
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if recoverFunc != nil {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						recoverFunc(ctx, v)
					} else {
						recoverFunc(ctx, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception))
					}
				}
			}
		}()
		goroutineFunc(ctx)
	}()
}

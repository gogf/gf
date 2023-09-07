// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpool

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Add pushes a new job to the pool.
// The job will be executed asynchronously.
func (p *Pool) Add(ctx context.Context, f Func) error {
	for p.closed.Val() {
		return gerror.NewCode(
			gcode.CodeInvalidOperation,
			"goroutine defaultPool is already closed",
		)
	}
	p.list.PushFront(&localPoolItem{
		Ctx:  ctx,
		Func: f,
	})
	// Check and fork new worker.
	p.checkAndForkNewGoroutineWorker()
	return nil
}

// AddWithRecover pushes a new job to the pool with specified recover function.
//
// The optional `recoverFunc` is called when any panic during executing of `userFunc`.
// If `recoverFunc` is not passed or given nil, it ignores the panic from `userFunc`.
// The job will be executed asynchronously.
func (p *Pool) AddWithRecover(ctx context.Context, userFunc Func, recoverFunc RecoverFunc) error {
	return p.Add(ctx, func(ctx context.Context) {
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
		userFunc(ctx)
	})
}

// Cap returns the capacity of the pool.
// This capacity is defined when pool is created.
// It returns -1 if there's no limit.
func (p *Pool) Cap() int {
	return p.limit
}

// Size returns current goroutine count of the pool.
func (p *Pool) Size() int {
	return p.count.Val()
}

// Jobs returns current job count of the pool.
// Note that, it does not return worker/goroutine count but the job/task count.
func (p *Pool) Jobs() int {
	return p.list.Size()
}

// IsClosed returns if pool is closed.
func (p *Pool) IsClosed() bool {
	return p.closed.Val()
}

// Close closes the goroutine pool, which makes all goroutines exit.
func (p *Pool) Close() {
	p.closed.Set(true)
}

// checkAndForkNewGoroutineWorker checks and creates a new goroutine worker.
// Note that the worker dies if the job function panics and the job has no recover handling.
func (p *Pool) checkAndForkNewGoroutineWorker() {
	// Check whether fork new goroutine or not.
	var n int
	for {
		n = p.count.Val()
		if p.limit != -1 && n >= p.limit {
			// No need fork new goroutine.
			return
		}
		if p.count.Cas(n, n+1) {
			// Use CAS to guarantee atomicity.
			break
		}
	}

	// Create job function in goroutine.
	go func() {
		defer p.count.Add(-1)

		var (
			listItem interface{}
			poolItem *localPoolItem
		)
		// Harding working, one by one, job never empty, worker never die.
		for !p.closed.Val() {
			listItem = p.list.PopBack()
			if listItem == nil {
				return
			}
			poolItem = listItem.(*localPoolItem)
			poolItem.Func(poolItem.Ctx)
		}
	}()
}

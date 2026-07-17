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

// Cap can change the capacity and returns the capacity of the pool before changed.
// This capacity is defined when pool is created. Pass newCap will change it.
// It returns -1 if there's no limit.
func (p *Pool) Cap(newCap ...int) int {
	if len(newCap) > 0 {
		cap := int64(newCap[0])
		if cap <= 0 {
			cap = -1
		}
		return int(p.limit.Swap(cap))
	} else {
		return int(p.limit.Load())
	}
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

// ClearJobs clear current all jobs and return how many cleared.
func (p *Pool) ClearJobs() (count int) {
	items := p.list.PopBackAll()
	return len(items)
}

// Pause pauses pool work. Jobs are kept in the queue and will be processed when the pool resumes.
func (p *Pool) Pause() bool {
	if p.IsClosed() {
		return false
	}
	if !p.paused.Swap(true) {
		if p.timer != nil {
			p.timer.Stop()
		}
	}
	return true
}

// IsPaused returns whether the pool is paused.
func (p *Pool) IsPaused() bool {
	return p.paused.Load()
}

// Resume resumes pool work.
func (p *Pool) Resume() bool {
	if p.IsClosed() {
		return false
	}
	if p.paused.Swap(false) {
		if p.timer != nil {
			p.timer.Start()
		}
	}
	return true
}

// IsClosed returns if pool is closed.
func (p *Pool) IsClosed() bool {
	return p.closed.Val()
}

// Close closes the goroutine pool, which makes all goroutines exit.
func (p *Pool) Close() {
	if p.closed.Cas(false, true) {
		if p.timer != nil {
			p.timer.Close()
			p.timer = nil
		}
	}
}

// checkAndForkNewGoroutineWorker checks and creates a new goroutine worker.
// Note that the worker dies if the job function panics and the job has no recover handling.
func (p *Pool) checkAndForkNewGoroutineWorker() {
	// Check whether fork new goroutine or not.
	if p.paused.Load() {
		return
	}
	var n int
	for {
		n = p.count.Val()
		if limit := p.limit.Load(); limit != -1 && int64(n) >= limit {
			// No need fork new goroutine.
			return
		}
		if p.count.Cas(n, n+1) {
			// Use CAS to guarantee atomicity.
			break
		}
	}

	// Create job function in goroutine.
	go p.asynchronousWorker()
}

func (p *Pool) asynchronousWorker() {
	var (
		n      int
		addVal = -1
	)
	defer func() { p.count.Add(addVal) }()
	// Harding working, one by one, job never empty, worker never die.
	for !p.closed.Val() && !p.paused.Load() {
		listItem := p.list.PopBack()
		if listItem == nil {
			return
		}
		listItem.Func(listItem.Ctx)
		// Check whether need reduce worker.
		n = p.count.Val()
		if limit := p.limit.Load(); limit > 0 && int64(n) > limit && p.count.Cas(n, n-1) {
			addVal = 0
			return
		}
	}
}

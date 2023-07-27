// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
package grpool

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/grand"
)

// Func is the pool function which contains context parameter.
type Func func(ctx context.Context)

// RecoverFunc is the pool runtime panic recover function which contains context parameter.
type RecoverFunc func(ctx context.Context, err error)

const (
	Closed  = 0 // Pool is Closed.
	Running = 1 // Pool is Running.
	Stopped = 2 // Pool is Stopped.
)

// Pool manages the goroutines using pool.
type Pool struct {
	limit int         // Max goroutine count limit.
	count *gtype.Int  // Current Running goroutine count.
	list  *glist.List // List for asynchronous job adding purpose.
	state *gtype.Int  // Pool state, 0: Closed, 1: Running. 2. Stopped
}

type localPoolItem struct {
	Ctx  context.Context
	Func Func
}

const (
	minTimerDuration = 500 * time.Millisecond
	maxTimerDuration = 1500 * time.Millisecond
)

// Default goroutine pool.
var (
	pool = New()
)

// New creates and returns a new goroutine pool object.
// The parameter `limit` is used to limit the max goroutine count,
// which is not limited in default.
func New(limit ...int) *Pool {
	p := &Pool{
		limit: -1,
		count: gtype.NewInt(),
		list:  glist.New(true),
		state: gtype.NewInt(Running), // the default state is Running
	}
	if len(limit) > 0 && limit[0] > 0 {
		p.limit = limit[0]
	}
	timerDuration := grand.D(minTimerDuration, maxTimerDuration)
	gtimer.Add(context.Background(), timerDuration, p.supervisor)
	return p
}

// Add pushes a new job to the pool using default goroutine pool.
// The job will be executed asynchronously.
func Add(ctx context.Context, f Func) error {
	return pool.Add(ctx, f)
}

// AddWithRecover pushes a new job to the pool with specified recover function.
// The optional `recoverFunc` is called when any panic during executing of `userFunc`.
// If `recoverFunc` is not passed or given nil, it ignores the panic from `userFunc`.
// The job will be executed asynchronously.
func AddWithRecover(ctx context.Context, userFunc Func, recoverFunc RecoverFunc) error {
	return pool.AddWithRecover(ctx, userFunc, recoverFunc)
}

// Size returns current goroutine count of default goroutine pool.
func Size() int {
	return pool.Size()
}

// Jobs returns current job count of default goroutine pool.
func Jobs() int {
	return pool.Jobs()
}

// Add pushes a new job to the pool.
// The job will be executed asynchronously.
func (p *Pool) Add(ctx context.Context, f Func) error {
	for p.state.Val() == Closed {
		return gerror.NewCode(
			gcode.CodeInvalidOperation,
			"goroutine pool is already Closed",
		)
	}
	p.list.PushFront(&localPoolItem{
		Ctx:  ctx,
		Func: f,
	})
	// Check and fork new worker.
	p.checkAndFork()
	return nil
}

// checkAndFork checks and creates a new goroutine worker.
// Note that the worker dies if the job function panics and the job has no recover handling.
func (p *Pool) checkAndFork() {
	// Check whether fork new goroutine or not.
	var n int
	for {
		n = p.count.Val()
		if p.limit != -1 && n >= p.limit {
			// No need fork new goroutine.
			return
		}
		if p.state.Val() != Running {
			// Pool is not Running, no need fork new goroutine.
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
		for p.state.Val() == Running {
			listItem = p.list.PopBack()
			if listItem == nil {
				return
			}
			poolItem = listItem.(*localPoolItem)
			poolItem.Func(poolItem.Ctx)
		}
	}()
}

// AddWithRecover pushes a new job to the pool with specified recover function.
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
						recoverFunc(ctx, gerror.Newf(`%+v`, exception))
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

// IsClosed returns if pool is Closed.
func (p *Pool) IsClosed() bool {
	return p.state.Val() == Closed
}

// IsStopped returns if pool is Stopped.
func (p *Pool) IsStopped() bool {
	return p.state.Val() == Stopped
}

// Close closes the goroutine pool, which makes all goroutines exit.
func (p *Pool) Close() {
	p.state.Set(Closed)
}

// IsRunning returns if pool is Running.
func (p *Pool) IsRunning() bool {
	return p.state.Val() == Running
}

// State returns the current state of the pool.
func (p *Pool) State() int {
	return p.state.Val()
}

// Start Starts the goroutine pool. You can call Stop() to stop the pool.
// Note That. You cant start a Closed pool.
func (p *Pool) Start() {
	p.state.Set(Running)
	p.checkAndFork()
}

// Stop stops the goroutine pool. You can call Start() again to restart the pool.
func (p *Pool) Stop() {
	p.state.Set(Stopped)
}

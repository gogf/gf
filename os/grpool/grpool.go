// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
package grpool

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/grand"
)

// Func is the pool function which contains context parameter.
type Func func(ctx context.Context)

type LimitChangerFunc func(ctx context.Context, value *atomic.Int64) (changed bool)

// RecoverFunc is the pool runtime panic recover function which contains context parameter.
type RecoverFunc func(ctx context.Context, exception error)

// Pool manages the goroutines using pool.
type Pool struct {
	limit        atomic.Int64                 // Max goroutine count limit.
	count        *gtype.Int                   // Current running goroutine count.
	list         *glist.TList[*localPoolItem] // List for asynchronous job adding purpose.
	closed       *gtype.Bool                  // Is pool closed or not.
	limitChanger LimitChangerFunc             // Function used to change max goroutine count limit. Let it nil to disable.
	parsed       atomic.Bool                  // Pool parsed to work
	timer        *gtimer.Entry
}

// PoolOption used to pass options
type PoolOption struct {
	Limit        int              // Max goroutine count limit.
	LimitChanger LimitChangerFunc // Function used to change max goroutine count limit. Let it nil to disable.
}

// localPoolItem is the job item storing in job list.
type localPoolItem struct {
	Ctx  context.Context // Context.
	Func Func            // Job function.
}

const (
	minSupervisorTimerDuration = 500 * time.Millisecond
	maxSupervisorTimerDuration = 1500 * time.Millisecond
)

// Default goroutine pool.
var (
	defaultPool = New()
)

// New creates and returns a new goroutine pool object.
// The parameter `limit` is used to limit the max goroutine count,
// which is not limited in default.
func New(limit ...int) *Pool {
	if len(limit) == 0 {
		return NewWithOption(PoolOption{
			Limit: -1,
		})
	} else {
		return NewWithOption(PoolOption{
			Limit: limit[0],
		})
	}
}

// New creates and returns a new goroutine pool object.
// The parameter `option` is used to limit the max goroutine count or set limit changer,
// which is not limited in default.
func NewWithOption(option ...PoolOption) *Pool {
	var o *PoolOption
	if len(option) > 0 {
		o = &option[0]
	} else {
		o = &PoolOption{}
	}
	var (
		pool = &Pool{
			limitChanger: nil,
			count:        gtype.NewInt(),
			list:         glist.NewT[*localPoolItem](true),
			closed:       gtype.NewBool(),
		}
		timerDuration = grand.D(
			minSupervisorTimerDuration,
			maxSupervisorTimerDuration,
		)
	)
	if o.Limit > 0 {
		pool.limit.Store(int64(o.Limit))
	} else {
		pool.limit.Store(-1)
	}
	if o.LimitChanger != nil {
		pool.limitChanger = o.LimitChanger
	}
	pool.timer = gtimer.AddSingleton(context.Background(), timerDuration, pool.supervisor)
	return pool
}

// Add pushes a new job to the default goroutine pool.
// The job will be executed asynchronously.
func Add(ctx context.Context, f Func) error {
	return defaultPool.Add(ctx, f)
}

// AddWithRecover pushes a new job to the default pool with specified recover function.
//
// The optional `recoverFunc` is called when any panic during executing of `userFunc`.
// If `recoverFunc` is not passed or given nil, it ignores the panic from `userFunc`.
// The job will be executed asynchronously.
func AddWithRecover(ctx context.Context, userFunc Func, recoverFunc RecoverFunc) error {
	return defaultPool.AddWithRecover(ctx, userFunc, recoverFunc)
}

// Size returns current goroutine count of default goroutine pool.
func Size() int {
	return defaultPool.Size()
}

// Jobs returns current job count of default goroutine pool.
func Jobs() int {
	return defaultPool.Jobs()
}

// Parse parse pool work.
func Parse() {
	defaultPool.Parse()
}

// IsClosed returns if pool is parsed.
func IsParsed() bool {
	return defaultPool.IsParsed()
}

// Resume resume pool work.
func Resume() {
	defaultPool.Resume()
}

// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
package grpool

import (
	"errors"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gtype"
)

// Goroutine Pool
type Pool struct {
	limit  int         // Max goroutine count limit.
	count  *gtype.Int  // Current running goroutine count.
	list   *glist.List // Job list for asynchronous job adding purpose.
	closed *gtype.Bool // Is pool closed or not.
}

// Default goroutine pool.
var pool = New()

// New creates and returns a new goroutine pool object.
// The parameter <limit> is used to limit the max goroutine count,
// which is not limited in default.
func New(limit ...int) *Pool {
	p := &Pool{
		limit:  -1,
		count:  gtype.NewInt(),
		list:   glist.New(true),
		closed: gtype.NewBool(),
	}
	if len(limit) > 0 && limit[0] > 0 {
		p.limit = limit[0]
	}
	return p
}

// Add pushes a new job to the pool using default goroutine pool.
// The job will be executed asynchronously.
func Add(f func()) error {
	return pool.Add(f)
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
func (p *Pool) Add(f func()) error {
	for p.closed.Val() {
		return errors.New("pool closed")
	}
	p.list.PushFront(f)
	// Check whether fork new goroutine or not.
	var n int
	for {
		n = p.count.Val()
		if p.limit != -1 && n >= p.limit {
			// No need fork new goroutine.
			return nil
		}
		if p.count.Cas(n, n+1) {
			// Use CAS to guarantee atomicity.
			break
		}
	}
	p.fork()
	return nil
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

// fork creates a new goroutine worker.
// Note that the worker dies if the job function panics.
func (p *Pool) fork() {
	go func() {
		defer p.count.Add(-1)

		var job interface{}
		for !p.closed.Val() {
			if job = p.list.PopBack(); job != nil {
				job.(func())()
			} else {
				return
			}
		}
	}()
}

// IsClosed returns if pool is closed.
func (p *Pool) IsClosed() bool {
	return p.closed.Val()
}

// Close closes the goroutine pool, which makes all goroutines exit.
func (p *Pool) Close() {
	p.closed.Set(true)
}

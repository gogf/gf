// Copyright 2017-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
package grpool

import (
	"github.com/gogf/gf/g/container/glist"
	"github.com/gogf/gf/g/container/gtype"
)

// Goroutine Pool
type Pool struct {
    limit  int           // Max goroutine count limit.
    count  *gtype.Int    // Current running goroutine count.
    list   *glist.List   // Job list for asynchronous job adding purpose.
    closed *gtype.Bool   // Is pool closed or not.
}

// Default goroutine pool.
var pool = New()

// New creates and returns a new goroutine pool object.
// The param <limit> is used to limit the max goroutine count,
// which is not limited in default.
func New(limit...int) *Pool {
    p := &Pool {
	    limit  : -1,
        count  : gtype.NewInt(),
        list   : glist.New(),
        closed : gtype.NewBool(),
    }
    if len(limit) > 0 && limit[0] > 0 {
    	p.limit = limit[0]
    }
    return p
}

// Add pushes a new job to the pool using default goroutine pool.
// The job will be executed asynchronously.
func Add(f func()) {
	pool.Add(f)
}

// Size returns current goroutine count of default goroutine pool.
func Size() int {
    return pool.count.Val()
}

// Jobs returns current job count of default goroutine pool.
func Jobs() int {
    return pool.list.Len()
}

// Add pushes a new job to the pool.
// The job will be executed asynchronously.
func (p *Pool) Add(f func()) {
    p.list.PushFront(f)
    // check whether to create a new goroutine or not.
    if p.count.Val() == p.limit {
		return
    }
	// ensure atomicity.
	if p.limit != -1 && p.count.Add(1) > p.limit {
		p.count.Add(-1)
		return
	}
    // fork a new goroutine to consume the job list.
	p.fork()
}


// Size returns current goroutine count of the pool.
func (p *Pool) Size() int {
    return p.count.Val()
}

// Jobs returns current job count of the pool.
func (p *Pool) Jobs() int {
    return p.list.Size()
}

// fork creates a new goroutine pool.
func (p *Pool) fork() {
    go func() {
    	defer p.count.Add(-1)
    	job := (interface{})(nil)
        for !p.closed.Val() {
        	if job = p.list.PopBack(); job != nil {
		        job.(func())()
	        } else {
	        	return
	        }
        }
    }()
}

// Close closes the goroutine pool, which makes all goroutines exit.
func (p *Pool) Close() {
	p.closed.Set(true)
}
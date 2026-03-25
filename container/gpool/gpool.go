// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gpool provides object-reusable concurrent-safe pool.
package gpool

import (
	"time"
)

// Pool is an Object-Reusable Pool.
type Pool struct {
	*TPool[any]
}

// NewFunc Creation function for object.
type NewFunc = TPoolNewFunc[any]

// ExpireFunc Destruction function for object.
type ExpireFunc = TPoolExpireFunc[any]

// New creates and returns a new object pool.
// To ensure execution efficiency, the expiration time cannot be modified once it is set.
//
// Note the expiration logic:
// ttl = 0 : not expired;
// ttl < 0 : immediate expired after use;
// ttl > 0 : timeout expired;
func New(ttl time.Duration, newFunc NewFunc, expireFunc ...ExpireFunc) *Pool {
	return &Pool{
		TPool: NewTPool(ttl, newFunc, expireFunc...),
	}
}

// Put puts an item to pool.
func (p *Pool) Put(value any) error {
	return p.TPool.Put(value)
}

// MustPut puts an item to pool, it panics if any error occurs.
func (p *Pool) MustPut(value any) {
	p.TPool.MustPut(value)
}

// Clear clears pool, which means it will remove all items from pool.
func (p *Pool) Clear() {
	p.TPool.Clear()
}

// Get picks and returns an item from pool. If the pool is empty and NewFunc is defined,
// it creates and returns one from NewFunc.
func (p *Pool) Get() (any, error) {
	return p.TPool.Get()
}

// Size returns the count of available items of pool.
func (p *Pool) Size() int {
	return p.TPool.Size()
}

// Close closes the pool. If `p` has ExpireFunc,
// then it automatically closes all items using this function before it's closed.
// Commonly you do not need to call this function manually.
func (p *Pool) Close() {
	p.TPool.Close()
}

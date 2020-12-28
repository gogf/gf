// Copyright GoFrame gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gpool provides object-reusable concurrent-safe pool.
package gpool

import (
	"errors"
	"time"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
)

// Pool is an Object-Reusable Pool.
type Pool struct {
	list    *glist.List                 // Available/idle items list.
	closed  *gtype.Bool                 // Whether the pool is closed.
	TTL     time.Duration               // Time To Live for pool items.
	NewFunc func() (interface{}, error) // Callback function to create pool item.
	// ExpireFunc is the for expired items destruction.
	// This function needs to be defined when the pool items
	// need to perform additional destruction operations.
	// Eg: net.Conn, os.File, etc.
	ExpireFunc func(interface{})
}

// Pool item.
type poolItem struct {
	value    interface{} // Item value.
	expireAt int64       // Expire timestamp in milliseconds.
}

// Creation function for object.
type NewFunc func() (interface{}, error)

// Destruction function for object.
type ExpireFunc func(interface{})

// New creates and returns a new object pool.
// To ensure execution efficiency, the expiration time cannot be modified once it is set.
//
// Note the expiration logic:
// ttl = 0 : not expired;
// ttl < 0 : immediate expired after use;
// ttl > 0 : timeout expired;
func New(ttl time.Duration, newFunc NewFunc, expireFunc ...ExpireFunc) *Pool {
	r := &Pool{
		list:    glist.New(true),
		closed:  gtype.NewBool(),
		TTL:     ttl,
		NewFunc: newFunc,
	}
	if len(expireFunc) > 0 {
		r.ExpireFunc = expireFunc[0]
	}
	gtimer.AddSingleton(time.Second, r.checkExpireItems)
	return r
}

// Put puts an item to pool.
func (p *Pool) Put(value interface{}) error {
	if p.closed.Val() {
		return errors.New("pool is closed")
	}
	item := &poolItem{
		value: value,
	}
	if p.TTL == 0 {
		item.expireAt = 0
	} else {
		// As for Golang version < 1.13, there's no method Milliseconds for time.Duration.
		// So we need calculate the milliseconds using its nanoseconds value.
		item.expireAt = gtime.TimestampMilli() + p.TTL.Nanoseconds()/1000000
	}
	p.list.PushBack(item)
	return nil
}

// Clear clears pool, which means it will remove all items from pool.
func (p *Pool) Clear() {
	if p.ExpireFunc != nil {
		for {
			if r := p.list.PopFront(); r != nil {
				p.ExpireFunc(r.(*poolItem).value)
			} else {
				break
			}
		}
	} else {
		p.list.RemoveAll()
	}

}

// Get picks and returns an item from pool. If the pool is empty and NewFunc is defined,
// it creates and returns one from NewFunc.
func (p *Pool) Get() (interface{}, error) {
	for !p.closed.Val() {
		if r := p.list.PopFront(); r != nil {
			f := r.(*poolItem)
			if f.expireAt == 0 || f.expireAt > gtime.TimestampMilli() {
				return f.value, nil
			} else if p.ExpireFunc != nil {
				// TODO: move expire function calling asynchronously from `Get` operation.
				p.ExpireFunc(f.value)
			}
		} else {
			break
		}
	}
	if p.NewFunc != nil {
		return p.NewFunc()
	}
	return nil, errors.New("pool is empty")
}

// Size returns the count of available items of pool.
func (p *Pool) Size() int {
	return p.list.Len()
}

// Close closes the pool. If <p> has ExpireFunc,
// then it automatically closes all items using this function before it's closed.
// Commonly you do not need call this function manually.
func (p *Pool) Close() {
	p.closed.Set(true)
}

// checkExpire removes expired items from pool in every second.
func (p *Pool) checkExpireItems() {
	if p.closed.Val() {
		// If p has ExpireFunc,
		// then it must close all items using this function.
		if p.ExpireFunc != nil {
			for {
				if r := p.list.PopFront(); r != nil {
					p.ExpireFunc(r.(*poolItem).value)
				} else {
					break
				}
			}
		}
		gtimer.Exit()
	}
	// All items do not expire.
	if p.TTL == 0 {
		return
	}
	// The latest item expire timestamp in milliseconds.
	var latestExpire int64 = -1
	// Retrieve the current timestamp in milliseconds, it expires the items
	// by comparing with this timestamp. It is not accurate comparison for
	// every items expired, but high performance.
	var timestampMilli = gtime.TimestampMilli()
	for {
		if latestExpire > timestampMilli {
			break
		}
		if r := p.list.PopFront(); r != nil {
			item := r.(*poolItem)
			latestExpire = item.expireAt
			// TODO improve the auto-expiration mechanism of the pool.
			if item.expireAt > timestampMilli {
				p.list.PushFront(item)
				break
			}
			if p.ExpireFunc != nil {
				p.ExpireFunc(item.value)
			}
		} else {
			break
		}
	}
}

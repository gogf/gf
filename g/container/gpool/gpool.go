// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gpool provides object-reusable concurrent-safe pool.
package gpool

import (
    "errors"
    "github.com/gogf/gf/g/container/glist"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/os/gtimer"
    "time"
)

// Object-Reusable Pool.
type Pool struct {
    list       *glist.List                // Available/idle list.
    closed     *gtype.Bool                // Whether the pool is closed.
    Expire     int64                      // Max idle time(ms), after which it is recycled.
    NewFunc    func()(interface{}, error) // Callback function to create item.
    ExpireFunc func(interface{})          // Expired destruction function for objects.
                                          // This function needs to be defined when the pool object
                                          // needs to perform additional destruction operations.
                                          // Eg: net.Conn, os.File, etc.
}

// Pool item.
type poolItem struct {
    expire int64               // Expire time(millisecond).
    value  interface{}         // Value.
}

// Creation function for object.
type NewFunc    func() (interface{}, error)

// Destruction function for object.
type ExpireFunc func(interface{})


// New returns a new object pool.
// To ensure execution efficiency, the expiration time cannot be modified once it is set.
// Expire:
// expire = 0 : not expired;
// expire < 0 : immediate recovery after use;
// expire > 0 : timeout recovery;
// Note that the expiration time unit is ** milliseconds **.
func New(expire int, newFunc NewFunc, expireFunc...ExpireFunc) *Pool {
    r := &Pool {
        list    : glist.New(),
        closed  : gtype.NewBool(),
        Expire  : int64(expire),
        NewFunc : newFunc,
    }
    if len(expireFunc) > 0 {
        r.ExpireFunc = expireFunc[0]
    }
    gtimer.AddSingleton(time.Second, r.checkExpire)
    return r
}

// Put puts an item to pool.
func (p *Pool) Put(value interface{}) {
    item := &poolItem {
        value : value,
    }
    if p.Expire == 0 {
        item.expire = 0
    } else {
        item.expire = gtime.Millisecond() + p.Expire
    }
    p.list.PushBack(item)
}

// Clear clears pool, which means it will remove all items from pool.
func (p *Pool) Clear() {
    p.list.RemoveAll()
}

// Get picks an item from pool.
func (p *Pool) Get() (interface{}, error) {
    for !p.closed.Val() {
        if r := p.list.PopFront(); r != nil {
            f := r.(*poolItem)
            if f.expire == 0 || f.expire > gtime.Millisecond() {
                return f.value, nil
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

// Close closes the pool.
func (p *Pool) Close() {
    p.closed.Set(true)
}

// checkExpire secondly removes expired items from pool.
func (p *Pool) checkExpire() {
    if p.closed.Val() {
        gtimer.Exit()
    }
    for {
        if r := p.list.PopFront(); r != nil {
            item := r.(*poolItem)
            if item.expire == 0 || item.expire > gtime.Millisecond() {
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
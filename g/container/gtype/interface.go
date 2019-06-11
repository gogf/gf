<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master

package gtype

import (
<<<<<<< HEAD
    "sync"
)

// 比较通用的并发安全数据类型
type Interface struct {
    mu  sync.RWMutex
    val interface{}
}

func NewInterface(value...interface{}) *Interface {
    if len(value) > 0 {
        return &Interface{val:value[0]}
    }
    return &Interface{}
}

func (t *Interface)Set(value interface{}) {
    t.mu.Lock()
    t.val = value
    t.mu.Unlock()
}

func (t *Interface)Val() interface{} {
    t.mu.RLock()
    b := t.val
    t.mu.RUnlock()
    return b
}

// 使用自定义方法执行加锁修改操作
func (t *Interface) LockFunc(f func(value interface{}) interface{}) {
    t.mu.Lock()
    t.val = f(t.val)
    t.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (t *Interface) RLockFunc(f func(value interface{})) {
    t.mu.RLock()
    f(t.val)
    t.mu.RUnlock()
=======
    "sync/atomic"
)

type Interface struct {
	value atomic.Value
}

// NewInterface returns a concurrent-safe object for interface{} type,
// with given initial value <value>.
func NewInterface(value...interface{}) *Interface {
    t := &Interface{}
    if len(value) > 0 && value[0] != nil {
        t.value.Store(value[0])
    }
    return t
}

// Clone clones and returns a new concurrent-safe object for interface{} type.
func (t *Interface) Clone() *Interface {
    return NewInterface(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
// Note: The parameter <value> cannot be nil.
func (t *Interface) Set(value interface{}) (old interface{}) {
    old = t.Val()
    t.value.Store(value)
    return
}

// Val atomically loads t.value.
func (t *Interface) Val() interface{} {
    return t.value.Load()
>>>>>>> upstream/master
}
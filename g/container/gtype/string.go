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

type String struct {
    mu  sync.RWMutex
    val string
}

func NewString(value...string) *String {
    if len(value) > 0 {
        return &String{val:value[0]}
    }
    return &String{}
}

func (t *String)Set(value string) {
    t.mu.Lock()
    t.val = value
    t.mu.Unlock()
}

func (t *String)Val() string {
    t.mu.RLock()
    s := t.val
    t.mu.RUnlock()
    return s
}

// 使用自定义方法执行加锁修改操作
func (t *String) LockFunc(f func(value string) string) {
    t.mu.Lock()
    t.val = f(t.val)
    t.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (t *String) RLockFunc(f func(value string)) {
    t.mu.RLock()
    f(t.val)
    t.mu.RUnlock()
=======
    "sync/atomic"
)

type String struct {
	value atomic.Value
}

// NewString returns a concurrent-safe object for string type,
// with given initial value <value>.
func NewString(value...string) *String {
    t := &String{}
    if len(value) > 0 {
        t.value.Store(value[0])
    }
    return t
}

// Clone clones and returns a new concurrent-safe object for string type.
func (t *String) Clone() *String {
    return NewString(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *String) Set(value string) (old string) {
    old = t.Val()
    t.value.Store(value)
    return
}

// Val atomically loads t.value.
func (t *String) Val() string {
    s := t.value.Load()
    if s != nil {
        return s.(string)
    }
    return ""
>>>>>>> upstream/master
}



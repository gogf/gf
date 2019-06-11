<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync"
)

type Bytes struct {
    mu  sync.RWMutex
    val []byte
}

func NewBytes(value...[]byte) *Bytes {
    if len(value) > 0 {
        return &Bytes{val:value[0]}
    }
    return &Bytes{}
}

func (t *Bytes)Set(value []byte) {
    t.mu.Lock()
    t.val = value
    t.mu.Unlock()
}

func (t *Bytes)Val() []byte {
    t.mu.RLock()
    b := t.val
    t.mu.RUnlock()
    return b
}

// 使用自定义方法执行加锁修改操作
func (t *Bytes) LockFunc(f func(value []byte) []byte) {
    t.mu.Lock()
    t.val = f(t.val)
    t.mu.Unlock()
}

// 使用自定义方法执行加锁读取操作
func (t *Bytes) RLockFunc(f func(value []byte)) {
    t.mu.RLock()
    f(t.val)
    t.mu.RUnlock()
}
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import "sync/atomic"

type Bytes struct {
	value atomic.Value
}

// NewBytes returns a concurrent-safe object for []byte type,
// with given initial value <value>.
func NewBytes(value...[]byte) *Bytes {
    t := &Bytes{}
    if len(value) > 0 {
        t.value.Store(value[0])
    }
    return t
}

// Clone clones and returns a new concurrent-safe object for []byte type.
func (t *Bytes) Clone() *Bytes {
    return NewBytes(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
// Note: The parameter <value> cannot be nil.
func (t *Bytes) Set(value []byte) (old []byte) {
    old = t.Val()
    t.value.Store(value)
    return
}

// Val atomically loads t.value.
func (t *Bytes) Val() []byte {
    if s := t.value.Load(); s != nil {
        return s.([]byte)
    }
    return nil
}
>>>>>>> upstream/master

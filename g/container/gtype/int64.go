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
    "sync/atomic"
)

type Int64 struct {
<<<<<<< HEAD
    val int64
}

func NewInt64(value...int64) *Int64 {
    if len(value) > 0 {
        return &Int64{val:value[0]}
=======
	value int64
}

// NewInt64 returns a concurrent-safe object for int64 type,
// with given initial value <value>.
func NewInt64(value...int64) *Int64 {
    if len(value) > 0 {
        return &Int64{
	        value : value[0],
		}
>>>>>>> upstream/master
    }
    return &Int64{}
}

<<<<<<< HEAD
func (t *Int64)Set(value int64) {
    atomic.StoreInt64(&t.val, value)
}

func (t *Int64)Val() int64 {
    return atomic.LoadInt64(&t.val)
}

func (t *Int64)Add(delta int64) int64 {
    return atomic.AddInt64(&t.val, delta)
=======
// Clone clones and returns a new concurrent-safe object for int64 type.
func (t *Int64) Clone() *Int64 {
    return NewInt64(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Int64) Set(value int64) (old int64) {
    return atomic.SwapInt64(&t.value, value)
}

// Val atomically loads t.value.
func (t *Int64) Val() int64 {
    return atomic.LoadInt64(&t.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Int64) Add(delta int64) int64 {
    return atomic.AddInt64(&t.value, delta)
>>>>>>> upstream/master
}
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

type Int32 struct {
<<<<<<< HEAD
    val int32
}

func NewInt32(value...int32) *Int32 {
    if len(value) > 0 {
        return &Int32{val:value[0]}
=======
	value int32
}

// NewInt32 returns a concurrent-safe object for int32 type,
// with given initial value <value>.
func NewInt32(value...int32) *Int32 {
    if len(value) > 0 {
        return &Int32{
	        value : value[0],
		}
>>>>>>> upstream/master
    }
    return &Int32{}
}

<<<<<<< HEAD
func (t *Int32)Set(value int32) {
    atomic.StoreInt32(&t.val, value)
}

func (t *Int32)Val() int32 {
    return atomic.LoadInt32(&t.val)
}

func (t *Int32)Add(delta int32) int32 {
    return atomic.AddInt32(&t.val, delta)
=======
// Clone clones and returns a new concurrent-safe object for int32 type.
func (t *Int32) Clone() *Int32 {
    return NewInt32(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Int32) Set(value int32) (old int32) {
    return atomic.SwapInt32(&t.value, value)
}

// Val atomically loads t.value.
func (t *Int32) Val() int32 {
    return atomic.LoadInt32(&t.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Int32) Add(delta int32) (new int32) {
    return atomic.AddInt32(&t.value, delta)
>>>>>>> upstream/master
}
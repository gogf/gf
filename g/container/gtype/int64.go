// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
    "sync/atomic"
)

type Int64 struct {
	value int64
}

// NewInt64 returns a concurrent-safe object for int64 type,
// with given initial value <value>.
func NewInt64(value...int64) *Int64 {
    if len(value) > 0 {
        return &Int64{
	        value : value[0],
		}
    }
    return &Int64{}
}

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
}
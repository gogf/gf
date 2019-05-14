// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
    "sync/atomic"
)

type Int struct {
	value int64
}

// NewInt returns a concurrent-safe object for int type,
// with given initial value <value>.
func NewInt(value...int) *Int {
    if len(value) > 0 {
        return &Int{
	        value : int64(value[0]),
		}
    }
    return &Int{}
}

// Clone clones and returns a new concurrent-safe object for int type.
func (t *Int) Clone() *Int {
    return NewInt(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Int) Set(value int) (old int) {
    return int(atomic.SwapInt64(&t.value, int64(value)))
}

// Val atomically loads t.value.
func (t *Int) Val() int {
    return int(atomic.LoadInt64(&t.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Int) Add(delta int) (new int) {
    return int(atomic.AddInt64(&t.value, int64(delta)))
}
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
    "sync/atomic"
)

type Uint64 struct {
	value uint64
}

// NewUint64 returns a concurrent-safe object for uint64 type,
// with given initial value <value>.
func NewUint64(value...uint64) *Uint64 {
    if len(value) > 0 {
        return &Uint64{
	        value : value[0],
		}
    }
    return &Uint64{}
}

// Clone clones and returns a new concurrent-safe object for uint64 type.
func (t *Uint64) Clone() *Uint64 {
    return NewUint64(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Uint64) Set(value uint64) (old uint64) {
    return atomic.SwapUint64(&t.value, value)
}

// Val atomically loads t.value.
func (t *Uint64) Val() uint64 {
    return atomic.LoadUint64(&t.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Uint64) Add(delta uint64) (new uint64) {
    return atomic.AddUint64(&t.value, delta)
}
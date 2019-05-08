// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"math"
	"sync/atomic"
	"unsafe"
)

type Float64 struct {
	value uint64
}

// NewFloat64 returns a concurrent-safe object for float64 type,
// with given initial value <value>.
func NewFloat64(value...float64) *Float64 {
    if len(value) > 0 {
        return &Float64{
	        value : math.Float64bits(value[0]),
		}
    }
    return &Float64{}
}

// Clone clones and returns a new concurrent-safe object for float64 type.
func (t *Float64) Clone() *Float64 {
    return NewFloat64(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Float64) Set(value float64) (old float64) {
    return math.Float64frombits(atomic.SwapUint64(&t.value, math.Float64bits(value)))
}

// Val atomically loads t.value.
func (t *Float64) Val() float64 {
    return math.Float64frombits(atomic.LoadUint64(&t.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Float64) Add(delta float64) (new float64) {
	for {
		old := math.Float64frombits(t.value)
		new  = old + delta
		if atomic.CompareAndSwapUint64(
			(*uint64)(unsafe.Pointer(&t.value)),
			math.Float64bits(old),
			math.Float64bits(new),
		) {
			break
		}
	}
	return
}

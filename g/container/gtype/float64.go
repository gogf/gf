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
    "sync/atomic"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

type Float64 struct {
    val uint64
}

func NewFloat64(value...float64) *Float64 {
    if len(value) > 0 {
        return &Float64{ val : float64ToUint64InBits(value[0]) }
=======
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
>>>>>>> upstream/master
    }
    return &Float64{}
}

<<<<<<< HEAD
func (t *Float64)Set(value float64) {
    atomic.StoreUint64(&t.val, float64ToUint64InBits(value) )
}

func (t *Float64)Val() float64 {
    return uint64ToFloat64InBits(atomic.LoadUint64(&t.val))
}

func (t *Float64)Add(delta float64) float64 {
    return uint64ToFloat64InBits(atomic.AddUint64(&t.val, float64ToUint64InBits(delta)))
}

// 通过二进制的方式将float64转换为uint64(都是64bits)
func float64ToUint64InBits(value float64) uint64 {
    b := gbinary.Encode(value)
    i := gbinary.DecodeToUint64(b)
    return i
}

// 通过二进制的方式将uint64转换为float64(都是64bits)
func uint64ToFloat64InBits(value uint64) float64 {
    b := gbinary.Encode(value)
    f := gbinary.DecodeToFloat64(b)
    return f
}
=======
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
>>>>>>> upstream/master

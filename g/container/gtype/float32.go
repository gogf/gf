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

type Float32 struct {
    val uint32
}

func NewFloat32(value...float32) *Float32 {
    if len(value) > 0 {
        return &Float32{ val : float32ToUint32InBits(value[0]) }
=======
	"math"
	"sync/atomic"
	"unsafe"
)

type Float32 struct {
	value uint32
}

// NewFloat32 returns a concurrent-safe object for float32 type,
// with given initial value <value>.
func NewFloat32(value...float32) *Float32 {
    if len(value) > 0 {
        return &Float32{
	        value : math.Float32bits(value[0]),
		}
>>>>>>> upstream/master
    }
    return &Float32{}
}

<<<<<<< HEAD
func (t *Float32)Set(value float32) {
    atomic.StoreUint32(&t.val, float32ToUint32InBits(value) )
}

func (t *Float32)Val() float32 {
    return uint32ToFloat32InBits(atomic.LoadUint32(&t.val))
}

func (t *Float32)Add(delta float32) float32 {
    return uint32ToFloat32InBits(atomic.AddUint32(&t.val, float32ToUint32InBits(delta)))
}

// 通过二进制的方式将float32转换为uint32(都是32bits)
func float32ToUint32InBits(value float32) uint32 {
    b := gbinary.Encode(value)
    i := gbinary.DecodeToUint32(b)
    return i
}

// 通过二进制的方式将uint32转换为float32(都是32bits)
func uint32ToFloat32InBits(value uint32) float32 {
    b := gbinary.Encode(value)
    f := gbinary.DecodeToFloat32(b)
    return f
=======
// Clone clones and returns a new concurrent-safe object for float32 type.
func (t *Float32) Clone() *Float32 {
    return NewFloat32(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Float32) Set(value float32) (old float32) {
    return math.Float32frombits(atomic.SwapUint32(&t.value, math.Float32bits(value)))
}

// Val atomically loads t.value.
func (t *Float32) Val() float32 {
    return math.Float32frombits(atomic.LoadUint32(&t.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Float32) Add(delta float32) (new float32) {
	for {
		old := math.Float32frombits(t.value)
		new  = old + delta
		if atomic.CompareAndSwapUint32(
			(*uint32)(unsafe.Pointer(&t.value)),
			math.Float32bits(old),
			math.Float32bits(new),
		) {
			break
		}
	}
	return
>>>>>>> upstream/master
}
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

type Float64 struct {
    val uint64
}

func NewFloat64(value...float64) *Float64 {
    if len(value) > 0 {
        return &Float64{ val : float64ToUint64InBits(value[0]) }
    }
    return &Float64{}
}

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
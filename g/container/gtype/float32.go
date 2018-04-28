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

type Float32 struct {
    val uint32
}

func NewFloat32(value...float32) *Float32 {
    if len(value) > 0 {
        return &Float32{ val : float32ToUint32InBits(value[0]) }
    }
    return &Float32{}
}

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
}
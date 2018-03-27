// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Float32 struct {
    val uint32
}

func NewFloat32(value...float32) *Float32 {
    if len(value) > 0 {
        return &Float32{val:uint32(value[0])}
    }
    return &Float32{}
}

func (t *Float32)Set(value float32) {
    atomic.StoreUint32(&t.val, uint32(value))
}

func (t *Float32)Get() int {
    return int(atomic.LoadUint32(&t.val))
}

func (t *Float32)Add(delta float32) int {
    return int(atomic.AddUint32(&t.val, uint32(delta)))
}
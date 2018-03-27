// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Float64 struct {
    val uint64
}

func NewFloat64(value...float32) *Float64 {
    if len(value) > 0 {
        return &Float64{val:uint64(value[0])}
    }
    return &Float64{}
}

func (t *Float64)Set(value float32) {
    atomic.StoreUint64(&t.val, uint64(value))
}

func (t *Float64)Val() int {
    return int(atomic.LoadUint64(&t.val))
}

func (t *Float64)Add(delta float32) int {
    return int(atomic.AddUint64(&t.val, uint64(delta)))
}
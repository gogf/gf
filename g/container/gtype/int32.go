// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Int32 struct {
    val int32
}

func NewInt32(value...int32) *Int32 {
    if len(value) > 0 {
        return &Int32{val:value[0]}
    }
    return &Int32{}
}

func (t *Int32)Set(value int32) {
    atomic.StoreInt32(&t.val, value)
}

func (t *Int32)Val() int32 {
    return atomic.LoadInt32(&t.val)
}

func (t *Int32)Add(delta int32) int32 {
    return atomic.AddInt32(&t.val, delta)
}
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Uint32 struct {
    val uint32
}

func NewUint32(value...uint32) *Uint32 {
    if len(value) > 0 {
        return &Uint32{val:value[0]}
    }
    return &Uint32{}
}

func (t *Uint32)Set(value uint32) {
    atomic.StoreUint32(&t.val, value)
}

func (t *Uint32)Val() uint32 {
    return atomic.LoadUint32(&t.val)
}

func (t *Uint32)Add(delta uint32) uint32 {
    return atomic.AddUint32(&t.val, delta)
}
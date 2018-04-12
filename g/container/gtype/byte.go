// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Byte struct {
    val int32
}

func NewByte(value...byte) *Byte {
    if len(value) > 0 {
        return &Byte{val : int32(value[0])}
    }
    return &Byte{}
}

func (t *Byte)Set(value byte) {
    atomic.StoreInt32(&t.val, int32(value))
}

func (t *Byte)Val() byte {
    return byte(atomic.LoadInt32(&t.val))
}

func (t *Byte)Add(delta int) byte {
    return byte(atomic.AddInt32(&t.val, int32(delta)))
}

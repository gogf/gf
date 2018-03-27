// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Uint struct {
    val uint64
}

func NewUint(value...uint) *Uint {
    if len(value) > 0 {
        return &Uint{val:uint64(value[0])}
    }
    return &Uint{}
}

func (t *Uint)Set(value uint) {
    atomic.StoreUint64(&t.val, uint64(value))
}

func (t *Uint)Val() uint {
    return uint(atomic.LoadUint64(&t.val))
}

func (t *Uint)Add(delta uint) int {
    return int(atomic.AddUint64(&t.val, uint64(delta)))
}
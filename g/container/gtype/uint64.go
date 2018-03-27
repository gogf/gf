// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Uint64 struct {
    val uint64
}

func NewUint64(value...uint64) *Uint64 {
    if len(value) > 0 {
        return &Uint64{val:value[0]}
    }
    return &Uint64{}
}

func (t *Uint64)Set(value uint64) {
    atomic.StoreUint64(&t.val, value)
}

func (t *Uint64)Val() uint64 {
    return atomic.LoadUint64(&t.val)
}

func (t *Uint64)Add(delta uint64) uint64 {
    return atomic.AddUint64(&t.val, delta)
}
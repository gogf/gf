// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import (
    "sync/atomic"
)

type Bytes struct {
    val atomic.Value
}

func NewBytes(value...[]byte) *Bytes {
    t := &Bytes{}
    if len(value) > 0 && value[0] != nil{
        t.val.Store(value[0])
    }
    return t
}

func (t *Bytes) Clone() *Bytes {
    return NewBytes(t.Val())
}

func (t *Bytes) Set(value []byte) {
    if value == nil {
        return
    }
    t.val.Store(value)
}

func (t *Bytes) Val() []byte {
    v := t.val.Load()
    if v != nil {
        return v.([]byte)
    }
    return nil
}

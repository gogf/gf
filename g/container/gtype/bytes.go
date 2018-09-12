// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtype

import "sync/atomic"

type Bytes struct {
    val atomic.Value
}

func NewBytes(value...[]byte) *Bytes {
    t := &Bytes{}
    if len(value) > 0 {
        t.val.Store(value[0])
    }
    return t
}

func (t *Bytes) Clone() *Bytes {
    return NewBytes(t.Val())
}

func (t *Bytes) Set(value []byte) {
    t.val.Store(value)
}

func (t *Bytes) Val() []byte {
    s := t.val.Load()
    if s != nil {
        return s.([]byte)
    }
    return nil
}

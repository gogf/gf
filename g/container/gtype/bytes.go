// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import "sync/atomic"

type Bytes struct {
	value atomic.Value
}

// NewBytes returns a concurrent-safe object for []byte type,
// with given initial value <value>.
func NewBytes(value...[]byte) *Bytes {
    t := &Bytes{}
    if len(value) > 0 {
        t.value.Store(value[0])
    }
    return t
}

// Clone clones and returns a new concurrent-safe object for []byte type.
func (t *Bytes) Clone() *Bytes {
    return NewBytes(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
// Note: The parameter <value> cannot be nil.
func (t *Bytes) Set(value []byte) (old []byte) {
    old = t.Val()
    t.value.Store(value)
    return
}

// Val atomically loads t.value.
func (t *Bytes) Val() []byte {
    if s := t.value.Load(); s != nil {
        return s.([]byte)
    }
    return nil
}

// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
    "sync/atomic"
)

type Byte struct {
	value int32
}

// NewByte returns a concurrent-safe object for byte type,
// with given initial value <value>.
func NewByte(value...byte) *Byte {
    if len(value) > 0 {
        return &Byte{
	        value : int32(value[0]),
		}
    }
    return &Byte{}
}

// Clone clones and returns a new concurrent-safe object for byte type.
func (t *Byte) Clone() *Byte {
    return NewByte(t.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (t *Byte) Set(value byte) (old byte) {
    return byte(atomic.SwapInt32(&t.value, int32(value)))
}

// Val atomically loads t.value.
func (t *Byte) Val() byte {
    return byte(atomic.LoadInt32(&t.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (t *Byte) Add(delta int) (new byte) {
    return byte(atomic.AddInt32(&t.value, int32(delta)))
}

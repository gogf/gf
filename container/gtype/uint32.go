// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"sync/atomic"
)

type Uint32 struct {
	value uint32
}

// NewUint32 returns a concurrent-safe object for uint32 type,
// with given initial value <value>.
func NewUint32(value ...uint32) *Uint32 {
	if len(value) > 0 {
		return &Uint32{
			value: value[0],
		}
	}
	return &Uint32{}
}

// Clone clones and returns a new concurrent-safe object for uint32 type.
func (v *Uint32) Clone() *Uint32 {
	return NewUint32(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Uint32) Set(value uint32) (old uint32) {
	return atomic.SwapUint32(&v.value, value)
}

// Val atomically loads t.value.
func (v *Uint32) Val() uint32 {
	return atomic.LoadUint32(&v.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Uint32) Add(delta uint32) (new uint32) {
	return atomic.AddUint32(&v.value, delta)
}

// Cas executes the compare-and-swap operation for value.
func (v *Uint32) Cas(old, new uint32) bool {
	return atomic.CompareAndSwapUint32(&v.value, old, new)
}

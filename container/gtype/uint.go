// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"github.com/gogf/gf/util/gconv"
	"strconv"
	"sync/atomic"
)

// Uint is a struct for concurrent-safe operation for type uint.
type Uint struct {
	value uint64
}

// NewUint creates and returns a concurrent-safe object for uint type,
// with given initial value <value>.
func NewUint(value ...uint) *Uint {
	if len(value) > 0 {
		return &Uint{
			value: uint64(value[0]),
		}
	}
	return &Uint{}
}

// Clone clones and returns a new concurrent-safe object for uint type.
func (v *Uint) Clone() *Uint {
	return NewUint(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Uint) Set(value uint) (old uint) {
	return uint(atomic.SwapUint64(&v.value, uint64(value)))
}

// Val atomically loads and returns t.value.
func (v *Uint) Val() uint {
	return uint(atomic.LoadUint64(&v.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Uint) Add(delta uint) (new uint) {
	return uint(atomic.AddUint64(&v.value, uint64(delta)))
}

// Cas executes the compare-and-swap operation for value.
func (v *Uint) Cas(old, new uint) (swapped bool) {
	return atomic.CompareAndSwapUint64(&v.value, uint64(old), uint64(new))
}

// String implements String interface for string printing.
func (v *Uint) String() string {
	return strconv.FormatUint(uint64(v.Val()), 10)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Uint) MarshalJSON() ([]byte, error) {
	return gconv.UnsafeStrToBytes(strconv.FormatUint(uint64(v.Val()), 10)), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Uint) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Uint(gconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Uint) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Uint(value))
	return nil
}

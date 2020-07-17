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

// Uint64 is a struct for concurrent-safe operation for type uint64.
type Uint64 struct {
	value uint64
}

// NewUint64 creates and returns a concurrent-safe object for uint64 type,
// with given initial value <value>.
func NewUint64(value ...uint64) *Uint64 {
	if len(value) > 0 {
		return &Uint64{
			value: value[0],
		}
	}
	return &Uint64{}
}

// Clone clones and returns a new concurrent-safe object for uint64 type.
func (v *Uint64) Clone() *Uint64 {
	return NewUint64(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Uint64) Set(value uint64) (old uint64) {
	return atomic.SwapUint64(&v.value, value)
}

// Val atomically loads and returns t.value.
func (v *Uint64) Val() uint64 {
	return atomic.LoadUint64(&v.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Uint64) Add(delta uint64) (new uint64) {
	return atomic.AddUint64(&v.value, delta)
}

// Cas executes the compare-and-swap operation for value.
func (v *Uint64) Cas(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&v.value, old, new)
}

// String implements String interface for string printing.
func (v *Uint64) String() string {
	return strconv.FormatUint(v.Val(), 10)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Uint64) MarshalJSON() ([]byte, error) {
	return gconv.UnsafeStrToBytes(strconv.FormatUint(v.Val(), 10)), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Uint64) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Uint64(gconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Uint64) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Uint64(value))
	return nil
}

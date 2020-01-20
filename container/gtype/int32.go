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

// Int32 is a struct for concurrent-safe operation for type int32.
type Int32 struct {
	value int32
}

// NewInt32 creates and returns a concurrent-safe object for int32 type,
// with given initial value <value>.
func NewInt32(value ...int32) *Int32 {
	if len(value) > 0 {
		return &Int32{
			value: value[0],
		}
	}
	return &Int32{}
}

// Clone clones and returns a new concurrent-safe object for int32 type.
func (v *Int32) Clone() *Int32 {
	return NewInt32(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Int32) Set(value int32) (old int32) {
	return atomic.SwapInt32(&v.value, value)
}

// Val atomically loads and returns t.value.
func (v *Int32) Val() int32 {
	return atomic.LoadInt32(&v.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Int32) Add(delta int32) (new int32) {
	return atomic.AddInt32(&v.value, delta)
}

// Cas executes the compare-and-swap operation for value.
func (v *Int32) Cas(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&v.value, old, new)
}

// String implements String interface for string printing.
func (v *Int32) String() string {
	return strconv.Itoa(int(v.Val()))
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Int32) MarshalJSON() ([]byte, error) {
	return gconv.UnsafeStrToBytes(strconv.Itoa(int(v.Val()))), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Int32) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Int32(gconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Int32) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Int32(value))
	return nil
}

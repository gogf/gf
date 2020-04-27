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

// Int64 is a struct for concurrent-safe operation for type int64.
type Int64 struct {
	value int64
}

// NewInt64 creates and returns a concurrent-safe object for int64 type,
// with given initial value <value>.
func NewInt64(value ...int64) *Int64 {
	if len(value) > 0 {
		return &Int64{
			value: value[0],
		}
	}
	return &Int64{}
}

// Clone clones and returns a new concurrent-safe object for int64 type.
func (v *Int64) Clone() *Int64 {
	return NewInt64(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Int64) Set(value int64) (old int64) {
	return atomic.SwapInt64(&v.value, value)
}

// Val atomically loads and returns t.value.
func (v *Int64) Val() int64 {
	return atomic.LoadInt64(&v.value)
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Int64) Add(delta int64) (new int64) {
	return atomic.AddInt64(&v.value, delta)
}

// Cas executes the compare-and-swap operation for value.
func (v *Int64) Cas(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&v.value, old, new)
}

// String implements String interface for string printing.
func (v *Int64) String() string {
	return strconv.FormatInt(v.Val(), 10)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Int64) MarshalJSON() ([]byte, error) {
	return gconv.UnsafeStrToBytes(strconv.FormatInt(v.Val(), 10)), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Int64) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Int64(gconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Int64) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Int64(value))
	return nil
}

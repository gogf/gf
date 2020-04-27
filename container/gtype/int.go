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

// Int is a struct for concurrent-safe operation for type int.
type Int struct {
	value int64
}

// NewInt creates and returns a concurrent-safe object for int type,
// with given initial value <value>.
func NewInt(value ...int) *Int {
	if len(value) > 0 {
		return &Int{
			value: int64(value[0]),
		}
	}
	return &Int{}
}

// Clone clones and returns a new concurrent-safe object for int type.
func (v *Int) Clone() *Int {
	return NewInt(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Int) Set(value int) (old int) {
	return int(atomic.SwapInt64(&v.value, int64(value)))
}

// Val atomically loads and returns t.value.
func (v *Int) Val() int {
	return int(atomic.LoadInt64(&v.value))
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Int) Add(delta int) (new int) {
	return int(atomic.AddInt64(&v.value, int64(delta)))
}

// Cas executes the compare-and-swap operation for value.
func (v *Int) Cas(old, new int) (swapped bool) {
	return atomic.CompareAndSwapInt64(&v.value, int64(old), int64(new))
}

// String implements String interface for string printing.
func (v *Int) String() string {
	return strconv.Itoa(v.Val())
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Int) MarshalJSON() ([]byte, error) {
	return gconv.UnsafeStrToBytes(strconv.Itoa(v.Val())), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Int) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Int(gconv.UnsafeBytesToStr(b)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Int) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Int(value))
	return nil
}

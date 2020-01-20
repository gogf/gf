// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"bytes"
	"github.com/gogf/gf/util/gconv"
	"sync/atomic"
)

// Bool is a struct for concurrent-safe operation for type bool.
type Bool struct {
	value int32
}

var (
	bytesTrue  = []byte("true")
	bytesFalse = []byte("false")
)

// NewBool creates and returns a concurrent-safe object for bool type,
// with given initial value <value>.
func NewBool(value ...bool) *Bool {
	t := &Bool{}
	if len(value) > 0 {
		if value[0] {
			t.value = 1
		} else {
			t.value = 0
		}
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for bool type.
func (v *Bool) Clone() *Bool {
	return NewBool(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Bool) Set(value bool) (old bool) {
	if value {
		old = atomic.SwapInt32(&v.value, 1) == 1
	} else {
		old = atomic.SwapInt32(&v.value, 0) == 1
	}
	return
}

// Val atomically loads and returns t.valueue.
func (v *Bool) Val() bool {
	return atomic.LoadInt32(&v.value) > 0
}

// Cas executes the compare-and-swap operation for value.
func (v *Bool) Cas(old, new bool) (swapped bool) {
	var oldInt32, newInt32 int32
	if old {
		oldInt32 = 1
	}
	if new {
		newInt32 = 1
	}
	return atomic.CompareAndSwapInt32(&v.value, oldInt32, newInt32)
}

// String implements String interface for string printing.
func (v *Bool) String() string {
	if v.Val() {
		return "true"
	}
	return "false"
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Bool) MarshalJSON() ([]byte, error) {
	if v.Val() {
		return bytesTrue, nil
	} else {
		return bytesFalse, nil
	}
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Bool) UnmarshalJSON(b []byte) error {
	v.Set(gconv.Bool(bytes.Trim(b, `"`)))
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Bool) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Bool(value))
	return nil
}

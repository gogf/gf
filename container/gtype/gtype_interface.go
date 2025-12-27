// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"sync/atomic"

	"github.com/gogf/gf/v2/internal/deepcopy"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

// Interface is a struct for concurrent-safe operation for type any.
type Interface struct {
	value atomic.Value
}

// NewInterface creates and returns a concurrent-safe object for any type,
// with given initial value `value`.
func NewInterface(value ...any) *Interface {
	t := &Interface{}
	if len(value) > 0 && value[0] != nil {
		t.value.Store(value[0])
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for any type.
func (v *Interface) Clone() *Interface {
	return NewInterface(v.Val())
}

// Set atomically stores `value` into t.value and returns the previous value of t.value.
// Note: The parameter `value` cannot be nil.
func (v *Interface) Set(value any) (old any) {
	old = v.Val()
	v.value.Store(value)
	return
}

// Val atomically loads and returns t.value.
func (v *Interface) Val() any {
	return v.value.Load()
}

// String implements String interface for string printing.
func (v *Interface) String() string {
	return gconv.String(v.Val())
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Interface) UnmarshalJSON(b []byte) error {
	var i any
	if err := json.UnmarshalUseNumber(b, &i); err != nil {
		return err
	}
	v.Set(i)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for `v`.
func (v *Interface) UnmarshalValue(value any) error {
	v.Set(value)
	return nil
}

// DeepCopy implements interface for deep copy of current type.
func (v *Interface) DeepCopy() any {
	if v == nil {
		return nil
	}
	return NewInterface(deepcopy.Copy(v.Val()))
}

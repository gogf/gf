// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"sync/atomic"
)

type String struct {
	value atomic.Value
}

// NewString returns a concurrent-safe object for string type,
// with given initial value <value>.
func NewString(value ...string) *String {
	t := &String{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

// Clone clones and returns a new concurrent-safe object for string type.
func (v *String) Clone() *String {
	return NewString(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *String) Set(value string) (old string) {
	old = v.Val()
	v.value.Store(value)
	return
}

// Val atomically loads t.value.
func (v *String) Val() string {
	s := v.value.Load()
	if s != nil {
		return s.(string)
	}
	return ""
}

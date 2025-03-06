// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/util/gconv"
)

// Val returns the current value of `v`.
func (v *Var) Val() any {
	if v == nil {
		return nil
	}
	if v.safe {
		if t, ok := v.value.(*gtype.Interface); ok {
			return t.Val()
		}
	}
	return v.value
}

// Interface is alias of Val.
func (v *Var) Interface() any {
	return v.Val()
}

// Bytes converts and returns `v` as []byte.
func (v *Var) Bytes() []byte {
	return gconv.Bytes(v.Val())
}

// String converts and returns `v` as string.
func (v *Var) String() string {
	return gconv.String(v.Val())
}

// Bool converts and returns `v` as bool.
func (v *Var) Bool() bool {
	return gconv.Bool(v.Val())
}

// Int converts and returns `v` as int.
func (v *Var) Int() int {
	return gconv.Int(v.Val())
}

// Int8 converts and returns `v` as int8.
func (v *Var) Int8() int8 {
	return gconv.Int8(v.Val())
}

// Int16 converts and returns `v` as int16.
func (v *Var) Int16() int16 {
	return gconv.Int16(v.Val())
}

// Int32 converts and returns `v` as int32.
func (v *Var) Int32() int32 {
	return gconv.Int32(v.Val())
}

// Int64 converts and returns `v` as int64.
func (v *Var) Int64() int64 {
	return gconv.Int64(v.Val())
}

// Uint converts and returns `v` as uint.
func (v *Var) Uint() uint {
	return gconv.Uint(v.Val())
}

// Uint8 converts and returns `v` as uint8.
func (v *Var) Uint8() uint8 {
	return gconv.Uint8(v.Val())
}

// Uint16 converts and returns `v` as uint16.
func (v *Var) Uint16() uint16 {
	return gconv.Uint16(v.Val())
}

// Uint32 converts and returns `v` as uint32.
func (v *Var) Uint32() uint32 {
	return gconv.Uint32(v.Val())
}

// Uint64 converts and returns `v` as uint64.
func (v *Var) Uint64() uint64 {
	return gconv.Uint64(v.Val())
}

// Float32 converts and returns `v` as float32.
func (v *Var) Float32() float32 {
	return gconv.Float32(v.Val())
}

// Float64 converts and returns `v` as float64.
func (v *Var) Float64() float64 {
	return gconv.Float64(v.Val())
}

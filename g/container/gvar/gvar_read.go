// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
    "github.com/gogf/gf/g/os/gtime"
    "time"
)

// 只读变量接口
type VarRead interface {
    Val() interface{}
    IsNil() bool
    Bytes() []byte
    String() string
    Bool() bool
    Int() int
    Int8() int8
    Int16() int16
    Int32() int32
    Int64() int64
    Uint() uint
    Uint8() uint8
    Uint16() uint16
    Uint32() uint32
    Uint64() uint64
    Float32() float32
    Float64() float64
    Interface() interface{}
    Ints() []int
    Floats() []float64
    Strings() []string
    Interfaces() []interface{}
    Time(format ...string) time.Time
    TimeDuration() time.Duration
    GTime(format...string) *gtime.Time
    Struct(objPointer interface{}, attrMapping ...map[string]string) error
}
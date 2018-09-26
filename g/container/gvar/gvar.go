// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 动态变量.
package gvar

import (
    "time"
    "gitee.com/johng/gf/g/util/gconv"
)

type Var struct {
    value interface{}
}

func New(value...interface{}) *Var {
    v := &Var{}
    if len(value) > 0 {
        v.value = value[0]
    }
    return v
}

func (v *Var) IsNil()          bool            { return v.value == nil }
func (v *Var) Bytes()          []byte          { return gconv.Bytes(v.value) }
func (v *Var) String()         string          { return gconv.String(v.value) }
func (v *Var) Bool()           bool            { return gconv.Bool(v.value) }

func (v *Var) Int()            int             { return gconv.Int(v.value) }
func (v *Var) Int8()           int8            { return gconv.Int8(v.value) }
func (v *Var) Int16()          int16           { return gconv.Int16(v.value) }
func (v *Var) Int32()          int32           { return gconv.Int32(v.value) }
func (v *Var) Int64()          int64           { return gconv.Int64(v.value) }

func (v *Var) Uint()           uint            { return gconv.Uint(v.value) }
func (v *Var) Uint8()          uint8           { return gconv.Uint8(v.value) }
func (v *Var) Uint16()         uint16          { return gconv.Uint16(v.value) }
func (v *Var) Uint32()         uint32          { return gconv.Uint32(v.value) }
func (v *Var) Uint64()         uint64          { return gconv.Uint64(v.value) }

func (v *Var) Float32()        float32         { return gconv.Float32(v.value) }
func (v *Var) Float64()        float64         { return gconv.Float64(v.value) }

func (v *Var) Strings()        []string        { return gconv.Strings(v.value) }

func (v *Var) Time(format...string) time.Time       { return gconv.Time(v.value, format...) }
func (v *Var) TimeDuration()        time.Duration   { return gconv.TimeDuration(v.value) }
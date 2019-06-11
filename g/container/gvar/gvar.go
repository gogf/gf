// Copyright 2018-2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gvar provides an universal variable type, like generics.
package gvar

import (
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/util/gconv"
    "time"
)

type Var struct {
    value interface{} // Underlying value.
    safe  bool        // Concurrent safe or not.
}

// New returns a new Var with given <value>.
// The param <unsafe> used to specify whether using Var in un-concurrent-safety,
// which is false in default, means concurrent-safe.
func New(value interface{}, unsafe...bool) *Var {
    v := &Var{}
    if len(unsafe) == 0 || !unsafe[0] {
        v.safe  = true
        v.value = gtype.NewInterface(value)
    } else {
        v.value = value
    }
    return v
}

// Set sets <value> to <v>, and returns the old value.
func (v *Var) Set(value interface{}) (old interface{}) {
    if v.safe {
        old = v.value.(*gtype.Interface).Set(value)
    } else {
        old = v.value
        v.value = value
    }
    return
}

// Val returns the current value of <v>.
func (v *Var) Val() interface{} {
    if v.safe {
        return v.value.(*gtype.Interface).Val()
    } else {
        return v.value
    }
}

// See Val().
func (v *Var) Interface() interface{} {
    return v.Val()
}

// Time converts and returns <v> as time.Time.
// The param <format> specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (v *Var) Time(format...string) time.Time {
    return gconv.Time(v.Val(), format...)
}

// TimeDuration converts and returns <v> as time.Duration.
// If value of <v> is string, then it uses time.ParseDuration for conversion.
func (v *Var) Duration() time.Duration {
    return gconv.Duration(v.Val())
}

// GTime converts and returns <v> as *gtime.Time.
// The param <format> specifies the format of the time string using gtime,
// eg: Y-m-d H:i:s.
func (v *Var) GTime(format...string) *gtime.Time {
    return gconv.GTime(v.Val(), format...)
}

// Struct maps value of <v> to <objPointer>.
// The param <objPointer> should be a pointer to a struct instance.
// The param <mapping> is used to specify the key-to-attribute mapping rules.
func (v *Var) Struct(pointer interface{}, mapping...map[string]string) error {
    return gconv.Struct(v.Val(), pointer, mapping...)
}

func (v *Var) IsNil()          bool            { return v.Val() == nil }
func (v *Var) Bytes()          []byte          { return gconv.Bytes(v.Val()) }
func (v *Var) String()         string          { return gconv.String(v.Val()) }
func (v *Var) Bool()           bool            { return gconv.Bool(v.Val()) }

func (v *Var) Int()            int             { return gconv.Int(v.Val()) }
func (v *Var) Int8()           int8            { return gconv.Int8(v.Val()) }
func (v *Var) Int16()          int16           { return gconv.Int16(v.Val()) }
func (v *Var) Int32()          int32           { return gconv.Int32(v.Val()) }
func (v *Var) Int64()          int64           { return gconv.Int64(v.Val()) }

func (v *Var) Uint()           uint            { return gconv.Uint(v.Val()) }
func (v *Var) Uint8()          uint8           { return gconv.Uint8(v.Val()) }
func (v *Var) Uint16()         uint16          { return gconv.Uint16(v.Val()) }
func (v *Var) Uint32()         uint32          { return gconv.Uint32(v.Val()) }
func (v *Var) Uint64()         uint64          { return gconv.Uint64(v.Val()) }

func (v *Var) Float32()        float32         { return gconv.Float32(v.Val()) }
func (v *Var) Float64()        float64         { return gconv.Float64(v.Val()) }

func (v *Var) Ints()           []int           { return gconv.Ints(v.Val()) }
func (v *Var) Floats()         []float64       { return gconv.Floats(v.Val()) }
func (v *Var) Strings()        []string        { return gconv.Strings(v.Val()) }
func (v *Var) Interfaces()     []interface{}   { return gconv.Interfaces(v.Val()) }

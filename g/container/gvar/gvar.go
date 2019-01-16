// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gvar provides an universal variable type, like generics.
//
// 通用动态变量.
package gvar

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gconv"
    "time"
)

type Var struct {
    value interface{} // 变量值
    safe  bool        // 当为true时,value为 *gtype.Interface 类型
}

// 创建一个动态变量，value参数可以为nil
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

// 创建一个只读动态变量，value参数可以为nil
func NewRead(value interface{}, unsafe...bool) VarRead {
    return VarRead(New(value, unsafe...))
}

// 返回动态变量的只读接口
func (v *Var) ReadOnly() VarRead {
    return VarRead(v)
}

func (v *Var) Set(value interface{}) (old interface{}) {
    if v.safe {
        old = v.value.(*gtype.Interface).Set(value)
    } else {
        old = v.value
        v.value = value
    }
    return
}

func (v *Var) Val() interface{} {
    if v.safe {
        return v.value.(*gtype.Interface).Val()
    } else {
        return v.value
    }
}

// Val() 别名
func (v *Var) Interface() interface{} {
    return v.Val()
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

func (v *Var) Time(format...string) time.Time {
    return gconv.Time(v.Val(), format...)
}
func (v *Var) TimeDuration() time.Duration {
    return gconv.TimeDuration(v.Val())
}

func (v *Var) GTime(format...string) *gtime.Time {
    return gconv.GTime(v.Val(), format...)
}

// 将变量转换为对象，注意 objPointer 参数必须为struct指针
func (v *Var) Struct(objPointer interface{}, attrMapping...map[string]string) error {
    return gconv.Struct(v.Val(), objPointer, attrMapping...)
}
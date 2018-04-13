// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 类型转换.
// 内部使用了bytes作为底层转换类型，效率很高。
package gconv

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

func Bytes(i interface{}) []byte {
    if i == nil {
        return nil
    }
    if r, ok := i.([]byte); ok {
        return r
    } else {
        return gbinary.Encode(i)
    }
}

// 基础的字符串类型转换
func String(i interface{}) string {
    if i == nil {
        return ""
    }
    if r, ok := i.(string); ok {
        return r
    } else {
        return string(Bytes(i))
    }
}

func Strings(i interface{}) []string {
    if i == nil {
        return nil
    }
    if r, ok := i.([]string); ok {
        return r
    } else if r, ok := i.([]interface{}); ok {
        strs := make([]string, len(r))
        for k, v := range r {
            strs[k] = String(v)
        }
        return strs
    }
    return []string{fmt.Sprintf("%v", i)}
}

//false: "", 0, false, off
func Bool(i interface{}) bool {
    if i == nil {
        return false
    }
    if v, ok := i.(bool); ok {
        return v
    }
    if s := String(i); s != "" && s != "0" && s != "false" && s != "off" {
        return true
    }
    return false
}

func Int(i interface{}) int {
    if i == nil {
        return 0
    }
    if v, ok := i.(int); ok {
        return v
    }
    return gbinary.DecodeToInt(Bytes(i))
}

func Int8(i interface{}) int8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int8); ok {
        return v
    }
    return gbinary.DecodeToInt8(Bytes(i))
}

func Int16(i interface{}) int16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int16); ok {
        return v
    }
    return gbinary.DecodeToInt16(Bytes(i))
}

func Int32(i interface{}) int32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int32); ok {
        return v
    }
    return gbinary.DecodeToInt32(Bytes(i))
}

func Int64(i interface{}) int64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int64); ok {
        return v
    }
    return gbinary.DecodeToInt64(Bytes(i))
}

func Uint(i interface{}) uint {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint); ok {
        return v
    }
    return gbinary.DecodeToUint(Bytes(i))
}

func Uint8(i interface{}) uint8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint8); ok {
        return v
    }
    return gbinary.DecodeToUint8(Bytes(i))
}

func Uint16(i interface{}) uint16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint16); ok {
        return v
    }
    return gbinary.DecodeToUint16(Bytes(i))
}

func Uint32(i interface{}) uint32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint32); ok {
        return v
    }
    return gbinary.DecodeToUint32(Bytes(i))
}

func Uint64(i interface{}) uint64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint64); ok {
        return v
    }
    return gbinary.DecodeToUint64(Bytes(i))
}

func Float32 (i interface{}) float32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float32); ok {
        return v
    }
    return gbinary.DecodeToFloat32(Bytes(i))
}

func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
    return gbinary.DecodeToFloat64(Bytes(i))
}

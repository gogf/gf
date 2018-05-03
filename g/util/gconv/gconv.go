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
    "time"
    "strconv"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

// 将变量i转换为字符串指定的类型t
func Convert(i interface{}, t string) interface{} {
    switch t {
        case "int":             return Int(i)
        case "int8":            return Int8(i)
        case "int16":           return Int16(i)
        case "int32":           return Int32(i)
        case "int64":           return Int64(i)
        case "uint":            return Uint(i)
        case "uint8":           return Uint8(i)
        case "uint16":          return Uint16(i)
        case "uint32":          return Uint32(i)
        case "uint64":          return Uint64(i)
        case "float32":         return Float32(i)
        case "float64":         return Float64(i)
        case "bool":            return Bool(i)
        case "string":          return String(i)
        case "[]byte":          return Bytes(i)
        case "time.Time":       return Time(i)
        case "time.Duration":   return TimeDuration(i)
        default:
            return i
    }
}

// 将变量i转换为time.Time类型
func Time(i interface{}) time.Time {
    s := String(i)
    t := int64(0)
    n := int64(0)
    if len(s) > 9 {
        t = Int64(s[0  : 10])
        if len(s) > 10 {
            n = Int64(s[11 : ])
        }
    }
    return time.Unix(t, n)
}

// 将变量i转换为time.Time类型
func TimeDuration(i interface{}) time.Duration {
    return time.Duration(Int64(i))
}

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
    switch value := i.(type) {
        case int:     return strconv.Itoa(value)
        case int8:    return strconv.Itoa(int(value))
        case int16:   return strconv.Itoa(int(value))
        case int32:   return strconv.Itoa(int(value))
        case int64:   return strconv.Itoa(int(value))
        case uint:    return strconv.FormatUint(uint64(value), 10)
        case uint8:   return strconv.FormatUint(uint64(value), 10)
        case uint16:  return strconv.FormatUint(uint64(value), 10)
        case uint32:  return strconv.FormatUint(uint64(value), 10)
        case uint64:  return strconv.FormatUint(uint64(value), 10)
        case float32: return strconv.FormatFloat(float64(value), 'f', -1, 64)
        case float64: return strconv.FormatFloat(value, 'f', -1, 64)
        case bool:    return strconv.FormatBool(value)
        case string:  return value
        case []byte:  return string(value)
        default:
            return fmt.Sprintf("%v", value)
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
    switch value := i.(type) {
        case int:     return value
        case int8:    return int(value)
        case int16:   return int(value)
        case int32:   return int(value)
        case int64:   return int(value)
        case uint:    return int(value)
        case uint8:   return int(value)
        case uint16:  return int(value)
        case uint32:  return int(value)
        case uint64:  return int(value)
        case float32: return int(value)
        case float64: return int(value)
        case bool:
            if value {
                return 1
            }
            return 0
        default:
            v, _ := strconv.Atoi(String(value))
            return v
    }
}

func Int8(i interface{}) int8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int8); ok {
        return v
    }
    return int8(Int(i))
}

func Int16(i interface{}) int16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int16); ok {
        return v
    }
    return int16(Int(i))
}

func Int32(i interface{}) int32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int32); ok {
        return v
    }
    return int32(Int(i))
}

func Int64(i interface{}) int64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int64); ok {
        return v
    }
    return int64(Int(i))
}

func Uint(i interface{}) uint {
    if i == nil {
        return 0
    }
    switch value := i.(type) {
        case int:     return uint(value)
        case int8:    return uint(value)
        case int16:   return uint(value)
        case int32:   return uint(value)
        case int64:   return uint(value)
        case uint:    return value
        case uint8:   return uint(value)
        case uint16:  return uint(value)
        case uint32:  return uint(value)
        case uint64:  return uint(value)
        case float32: return uint(value)
        case float64: return uint(value)
        case bool:
            if value {
                return 1
            }
            return 0
        default:
            v, _ := strconv.ParseUint(String(value), 10, 64)
            return uint(v)
    }
}

func Uint8(i interface{}) uint8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint8); ok {
        return v
    }
    return uint8(Uint(i))
}

func Uint16(i interface{}) uint16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint16); ok {
        return v
    }
    return uint16(Uint(i))
}

func Uint32(i interface{}) uint32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint32); ok {
        return v
    }
    return uint32(Uint(i))
}

func Uint64(i interface{}) uint64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint64); ok {
        return v
    }
    return uint64(Uint(i))
}

func Float32 (i interface{}) float32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float32); ok {
        return v
    }
    v, _ := strconv.ParseFloat(String(i), 32)
    return float32(v)
}

func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
    v, _ := strconv.ParseFloat(String(i), 64)
    return v
}



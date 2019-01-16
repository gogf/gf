// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gconv implements powerful and easy-to-use converting functionality for any types of variables.
// 
// 类型转换, 
// 内部使用了bytes作为底层转换类型，效率很高。
package gconv

import (
    "strconv"
    "encoding/json"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "strings"
)

// 转换为string类型的接口
type apiString interface {
    String() string
}

// 将变量i转换为字符串指定的类型t，非必须参数extraParams用以额外的参数传递
func Convert(i interface{}, t string, extraParams...interface{}) interface{} {
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
        case "[]int":           return Ints(i)
        case "[]string":        return Strings(i)
        case "time.Time":
            if len(extraParams) > 0 {
                return Time(i, String(extraParams[0]))
            }
            return Time(i)

        case "time.Duration":   return TimeDuration(i)
        default:
            return i
    }
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
        case float32: return strconv.FormatFloat(float64(value), 'f', -1, 32)
        case float64: return strconv.FormatFloat(value, 'f', -1, 64)
        case bool:    return strconv.FormatBool(value)
        case string:  return value
        case []byte:  return string(value)
        default:
            if f, ok := value.(apiString); ok {
                // 如果变量实现了String()接口，那么使用该接口执行转换
                return f.String()
            } else {
                // 默认使用json进行字符串转换
                jsonContent, _ := json.Marshal(value)
                return string(jsonContent)
            }
    }
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
            return int(Float64(value))
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
        case int:
            if value < 0 {
                value = -value
            }
            return uint(value)
        case int8:
            if value < 0 {
                value = -value
            }
            return uint(value)
        case int16:
            if value < 0 {
                value = -value
            }
            return uint(value)
        case int32:
            if value < 0 {
                value = -value
            }
            return uint(value)
        case int64:
            if value < 0 {
                value = -value
            }
            return uint(value)
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
            v := Float64(value)
            if v < 0 {
                v = -v
            }
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
    v, _ := strconv.ParseFloat(strings.TrimSpace(String(i)), 64)
    return float32(v)
}

func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
    v, _ := strconv.ParseFloat(strings.TrimSpace(String(i)), 64)
    return v
}


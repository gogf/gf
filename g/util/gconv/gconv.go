<<<<<<< HEAD
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
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gconv implements powerful and easy-to-use converting functionality for any types of variables.
package gconv

import (
    "encoding/json"
    "github.com/gogf/gf/g/encoding/gbinary"
    "reflect"
    "strconv"
    "strings"
)

// Type assert api for String().
type apiString interface {
    String() string
}

// Type assert api for Error().
type apiError interface {
	Error() string
}

var (
    // Empty strings.
    emptyStringMap = map[string]struct{}{
        ""      : struct {}{},
        "0"     : struct {}{},
        "off"   : struct {}{},
        "false" : struct {}{},
    }
)


// Convert converts the variable <i> to the type <t>, the type <t> is specified by string.
// The unnecessary parameter <params> is used for additional parameter passing.
func Convert(i interface{}, t string, params...interface{}) interface{} {
>>>>>>> upstream/master
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
<<<<<<< HEAD
        case "time.Time":       return Time(i)
        case "time.Duration":   return TimeDuration(i)
=======
        case "[]int":           return Ints(i)
        case "[]string":        return Strings(i)

        case "Time", "time.Time":
            if len(params) > 0 {
                return Time(i, String(params[0]))
            }
            return Time(i)

        case "gtime.Time":
            if len(params) > 0 {
                return GTime(i, String(params[0]))
            }
            return *GTime(i)

        case "GTime", "*gtime.Time":
            if len(params) > 0 {
                return GTime(i, String(params[0]))
            }
            return GTime(i)

        case "Duration", "time.Duration":
        	return Duration(i)
>>>>>>> upstream/master
        default:
            return i
    }
}

<<<<<<< HEAD
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

=======
// Byte converts <i> to byte.
func Byte(i interface{}) byte {
	if v, ok := i.(byte); ok {
		return v
	}
	return byte(Uint8(i))
}

// Bytes converts <i> to []byte.
>>>>>>> upstream/master
func Bytes(i interface{}) []byte {
    if i == nil {
        return nil
    }
<<<<<<< HEAD
    if r, ok := i.([]byte); ok {
        return r
    } else {
        return gbinary.Encode(i)
    }
}

// 基础的字符串类型转换
=======
    switch value := i.(type) {
        case string:  return []byte(value)
        case []byte:  return value
        default:
            return gbinary.Encode(i)
    }
}

// Rune converts <i> to rune.
func Rune(i interface{}) rune {
	if v, ok := i.(rune); ok {
		return v
	}
	return rune(Int32(i))
}

// Runes converts <i> to []rune.
func Runes(i interface{}) []rune {
	if v, ok := i.([]rune); ok {
		return v
	}
	return []rune(String(i))
}


// String converts <i> to string.
>>>>>>> upstream/master
func String(i interface{}) string {
    if i == nil {
        return ""
    }
    switch value := i.(type) {
<<<<<<< HEAD
        case int:     return strconv.Itoa(value)
        case int8:    return strconv.Itoa(int(value))
        case int16:   return strconv.Itoa(int(value))
        case int32:   return strconv.Itoa(int(value))
        case int64:   return strconv.Itoa(int(value))
=======
        case int:     return strconv.FormatInt(int64(value), 10)
        case int8:    return strconv.Itoa(int(value))
        case int16:   return strconv.Itoa(int(value))
        case int32:   return strconv.Itoa(int(value))
        case int64:   return strconv.FormatInt(int64(value), 10)
>>>>>>> upstream/master
        case uint:    return strconv.FormatUint(uint64(value), 10)
        case uint8:   return strconv.FormatUint(uint64(value), 10)
        case uint16:  return strconv.FormatUint(uint64(value), 10)
        case uint32:  return strconv.FormatUint(uint64(value), 10)
        case uint64:  return strconv.FormatUint(uint64(value), 10)
<<<<<<< HEAD
        case float32: return strconv.FormatFloat(float64(value), 'f', -1, 64)
=======
        case float32: return strconv.FormatFloat(float64(value), 'f', -1, 32)
>>>>>>> upstream/master
        case float64: return strconv.FormatFloat(value, 'f', -1, 64)
        case bool:    return strconv.FormatBool(value)
        case string:  return value
        case []byte:  return string(value)
<<<<<<< HEAD
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
=======
        case []rune:  return string(value)
        default:
            if f, ok := value.(apiString); ok {
                // If the variable implements the String() interface,
                // then use that interface to perform the conversion
                return f.String()
            } else if f, ok := value.(apiError); ok {
	            // If the variable implements the Error() interface,
	            // then use that interface to perform the conversion
	            return f.Error()
            } else {
                // Finally we use json.Marshal to convert.
                jsonContent, _ := json.Marshal(value)
                return string(jsonContent)
            }
    }
}

// Bool converts <i> to bool.
// It returns false if <i> is: false, "", 0, "false", "off", empty slice/map.
>>>>>>> upstream/master
func Bool(i interface{}) bool {
    if i == nil {
        return false
    }
    if v, ok := i.(bool); ok {
        return v
    }
<<<<<<< HEAD
    if s := String(i); s != "" && s != "0" && s != "false" && s != "off" {
        return true
    }
    return false
}

=======
    if s, ok := i.(string); ok {
        if _, ok := emptyStringMap[s]; ok {
            return false
        }
        return true
    }
    rv := reflect.ValueOf(i)
    switch rv.Kind() {
        case reflect.Ptr:    return !rv.IsNil()
        case reflect.Map:    fallthrough
        case reflect.Array:  fallthrough
        case reflect.Slice:  return rv.Len() != 0
        case reflect.Struct: return true
        default:
            s := String(i)
            if _, ok := emptyStringMap[s]; ok {
                return false
            }
            return true

    }
}

// Int converts <i> to int.
>>>>>>> upstream/master
func Int(i interface{}) int {
    if i == nil {
        return 0
    }
<<<<<<< HEAD
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

=======
    if v, ok := i.(int); ok {
        return v
    }
    return int(Int64(i))
}

// Int8 converts <i> to int8.
>>>>>>> upstream/master
func Int8(i interface{}) int8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int8); ok {
        return v
    }
<<<<<<< HEAD
    return int8(Int(i))
}

=======
    return int8(Int64(i))
}

// Int16 converts <i> to int16.
>>>>>>> upstream/master
func Int16(i interface{}) int16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int16); ok {
        return v
    }
<<<<<<< HEAD
    return int16(Int(i))
}

=======
    return int16(Int64(i))
}

// Int32 converts <i> to int32.
>>>>>>> upstream/master
func Int32(i interface{}) int32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int32); ok {
        return v
    }
<<<<<<< HEAD
    return int32(Int(i))
}

=======
    return int32(Int64(i))
}

// Int64 converts <i> to int64.
>>>>>>> upstream/master
func Int64(i interface{}) int64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(int64); ok {
        return v
    }
<<<<<<< HEAD
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
=======
    switch value := i.(type) {
        case int:     return int64(value)
        case int8:    return int64(value)
        case int16:   return int64(value)
        case int32:   return int64(value)
        case int64:   return value
        case uint:    return int64(value)
        case uint8:   return int64(value)
        case uint16:  return int64(value)
        case uint32:  return int64(value)
        case uint64:  return int64(value)
        case float32: return int64(value)
        case float64: return int64(value)
>>>>>>> upstream/master
        case bool:
            if value {
                return 1
            }
            return 0
        default:
<<<<<<< HEAD
            v, _ := strconv.ParseUint(String(value), 10, 64)
            return uint(v)
    }
}

=======
            s := String(value)
            // Hexadecimal
            if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
                if v, e := strconv.ParseInt(s[2 : ], 16, 64); e == nil {
                    return v
                }
            }
            // Octal
            if len(s) > 1 && s[0] == '0' {
                if v, e := strconv.ParseInt(s[1 : ], 8, 64); e == nil {
                    return v
                }
            }
            // Decimal
            if v, e := strconv.ParseInt(s, 10, 64); e == nil {
                return v
            }
            // Float64
            return int64(Float64(value))
    }
}

// Uint converts <i> to uint.
func Uint(i interface{}) uint {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint); ok {
        return v
    }
    return uint(Uint64(i))
}

// Uint8 converts <i> to uint8.
>>>>>>> upstream/master
func Uint8(i interface{}) uint8 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint8); ok {
        return v
    }
<<<<<<< HEAD
    return uint8(Uint(i))
}

=======
    return uint8(Uint64(i))
}

// Uint16 converts <i> to uint16.
>>>>>>> upstream/master
func Uint16(i interface{}) uint16 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint16); ok {
        return v
    }
<<<<<<< HEAD
    return uint16(Uint(i))
}

=======
    return uint16(Uint64(i))
}

// Uint32 converts <i> to uint32.
>>>>>>> upstream/master
func Uint32(i interface{}) uint32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint32); ok {
        return v
    }
<<<<<<< HEAD
    return uint32(Uint(i))
}

=======
    return uint32(Uint64(i))
}

// Uint64 converts <i> to uint64.
>>>>>>> upstream/master
func Uint64(i interface{}) uint64 {
    if i == nil {
        return 0
    }
<<<<<<< HEAD
    if v, ok := i.(uint64); ok {
        return v
    }
    return uint64(Uint(i))
}

=======
    switch value := i.(type) {
        case int:     return uint64(value)
        case int8:    return uint64(value)
        case int16:   return uint64(value)
        case int32:   return uint64(value)
        case int64:   return uint64(value)
        case uint:    return uint64(value)
        case uint8:   return uint64(value)
        case uint16:  return uint64(value)
        case uint32:  return uint64(value)
        case uint64:  return value
        case float32: return uint64(value)
        case float64: return uint64(value)
        case bool:
            if value {
                return 1
            }
            return 0
        default:
            s := String(value)
            // Hexadecimal
            if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
                if v, e := strconv.ParseUint(s[2 : ], 16, 64); e == nil {
                    return v
                }
            }
            // Octal
            if len(s) > 1 && s[0] == '0' {
                if v, e := strconv.ParseUint(s[1 : ], 8, 64); e == nil {
                    return v
                }
            }
            // Decimal
            if v, e := strconv.ParseUint(s, 10, 64); e == nil {
                return v
            }
            // Float64
            return uint64(Float64(value))
    }
}

// Float32 converts <i> to float32.
>>>>>>> upstream/master
func Float32 (i interface{}) float32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float32); ok {
        return v
    }
<<<<<<< HEAD
    v, _ := strconv.ParseFloat(String(i), 32)
    return float32(v)
}

=======
    v, _ := strconv.ParseFloat(strings.TrimSpace(String(i)), 64)
    return float32(v)
}

// Float64 converts <i> to float64.
>>>>>>> upstream/master
func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
<<<<<<< HEAD
    v, _ := strconv.ParseFloat(String(i), 64)
    return v
}


=======
    v, _ := strconv.ParseFloat(strings.TrimSpace(String(i)), 64)
    return v
}

>>>>>>> upstream/master

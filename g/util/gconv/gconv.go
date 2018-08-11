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
    "encoding/json"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gstr"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "github.com/fatih/structs"
    "strings"
    "reflect"
)

// 将变量i转换为字符串指定的类型t
func Convert(i interface{}, t string, params...interface{}) interface{} {
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
        case "time.Time":
            if len(params) > 0 {
                return Time(i, String(params[0]))
            }
            return Time(i)

        case "time.Duration":   return TimeDuration(i)
        default:
            return i
    }
}

// 将变量i转换为time.Time类型
func Time(i interface{}, format...string) time.Time {
    s := String(i)
    // 优先使用用户输入日期格式进行转换
    if len(format) > 0 {
        t, _ := gtime.StrToTimeFormat(s, format[0])
        return t
    }
    if gstr.IsNumeric(s) {
        return gtime.NewFromTimeStamp(Int64(s)).Time
    } else {
        t, _ := gtime.StrToTime(s)
        return t
    }
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
        case float32: return strconv.FormatFloat(float64(value), 'f', -1, 32)
        case float64: return strconv.FormatFloat(value, 'f', -1, 64)
        case bool:    return strconv.FormatBool(value)
        case string:  return value
        case []byte:  return string(value)
        default:
            // 默认使用json进行字符串转换
            jsonContent, _ := json.Marshal(value)
            return string(jsonContent)
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
    v, _ := strconv.ParseFloat(String(i), 64)
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

// 将params键值对参数映射到对应的struct对象属性上，第三个参数mapping为非必需，表示自定义名称与属性名称的映射关系。
// 需要注意：
// 1、第二个参数为struct对象指针；
// 2、struct对象的**公开属性(首字母大写)**才能被映射赋值；
// 3、map中的键名可以为小写，映射转换时会自动将键名首字母转为大写做匹配映射，如果无法匹配则忽略；
func MapToStruct(params map[string]interface{}, object interface{}, mapping...map[string]string) error {
    tagmap := make(map[string]string)
    fields := structs.Fields(object)
    // 将struct中定义的属性转换名称构建称tagmap
    for _, field := range fields {
        if tag := field.Tag("gconv"); tag != "" {
            for _, v := range strings.Split(tag, ",") {
                tagmap[strings.TrimSpace(v)] = field.Name()
            }
        }
    }
    elem := reflect.ValueOf(object).Elem()
    dmap := make(map[string]bool)
    // 首先按照传递的映射关系进行匹配
    if len(mapping) > 0 {
        for mappingk, mappingv := range mapping[0] {
            if v, ok := params[mappingk]; ok {
                dmap[mappingv] = true
                bindVarToStruct(elem, mappingv, v)
            }
        }
    }
    // 其次匹配对象定义时绑定的属性名称
    for tagk, tagv := range tagmap {
        if _, ok := dmap[tagv]; ok {
            continue
        }
        if v, ok := params[tagk]; ok {
            dmap[tagv] = true
            bindVarToStruct(elem, tagv, v)
        }
    }
    // 最后按照默认规则进行匹配
    for mapk, mapv := range params {
        name := gstr.UcFirst(mapk)
        if _, ok := dmap[name]; ok {
            continue
        }
        // 后续tag逻辑中会处理的key(重复的键名)这里便不处理
        if _, ok := tagmap[mapk]; !ok {
            bindVarToStruct(elem, name, mapv)
        }
    }
    return nil
}

// 将参数值绑定到对象指定名称的属性上
func bindVarToStruct(elem reflect.Value, name string, value interface{}) {
    structFieldValue := elem.FieldByName(name)
    // 键名与对象属性匹配检测
    if !structFieldValue.IsValid() {
        return
    }
    // CanSet的属性必须为公开属性(首字母大写)
    if !structFieldValue.CanSet() {
        return
    }
    // 必须将value转换为struct属性的数据类型，这里必须用到gconv包
    structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
}

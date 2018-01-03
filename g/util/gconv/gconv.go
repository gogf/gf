// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 类型转换.
// 如果给定的interface{}参数不是指定转换的输出类型，那么会进行强制转换，效率会比较低，
// 建议已知类型的转换自行调用相关方法来单独处理。
package gconv

import (
    "fmt"
    "strconv"
)

func Bytes(i interface{}) []byte {
    if i == nil {
        return nil
    }
    if r, ok := i.([]byte); ok {
        return r
    } else {
        return []byte(String(i))
    }
}

func String(i interface{}) string {
    if i == nil {
        return ""
    }
    if r, ok := i.(string); ok {
        return r
    } else {
        return fmt.Sprintf("%v", i)
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
    v, _ := strconv.Atoi(String(i))
    return v
}

func Uint (i interface{}) uint {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint); ok {
        return v
    }
    v, _ := strconv.ParseUint(String(i), 10, 8)
    return uint(v)
}

func Float32 (i interface{}) float32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float32); ok {
        return v
    }
    v, _ := strconv.ParseFloat(String(i), 8)
    return float32(v)
}

func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
    v, _ := strconv.ParseFloat(String(i), 8)
    return v
}

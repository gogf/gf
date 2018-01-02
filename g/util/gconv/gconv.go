// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 数据基本类型强制转换
package gconv

import (
    "fmt"
    "strconv"
)

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
    v, _ := strconv.Atoi(fmt.Sprintf("%v", i))
    return v
}

func Uint (i interface{}) uint {
    if i == nil {
        return 0
    }
    if v, ok := i.(uint); ok {
        return v
    }
    v, _ := strconv.ParseUint(fmt.Sprintf("%v", i), 10, 8)
    return uint(v)
}

func Float32 (i interface{}) float32 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float32); ok {
        return v
    }
    v, _ := strconv.ParseFloat(fmt.Sprintf("%v", i), 8)
    return float32(v)
}

func Float64 (i interface{}) float64 {
    if i == nil {
        return 0
    }
    if v, ok := i.(float64); ok {
        return v
    }
    v, _ := strconv.ParseFloat(fmt.Sprintf("%v", i), 8)
    return v
}

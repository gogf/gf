// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv

import (
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gstr"
)

// 将变量i转换为time.Time类型
func Time(i interface{}, format...string) time.Time {
    s := String(i)
    // 优先使用用户输入日期格式进行转换
    if len(format) > 0 {
        t, _ := gtime.StrToTimeFormat(s, format[0])
        return t.Time
    }
    if gstr.IsNumeric(s) {
        return gtime.NewFromTimeStamp(Int64(s)).Time
    } else {
        t, _ := gtime.StrToTime(s)
        return t.Time
    }
}

// 将变量i转换为time.Time类型
func TimeDuration(i interface{}) time.Duration {
    return time.Duration(Int64(i))
}
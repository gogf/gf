// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
    "time"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/text/gstr"
)

// 将变量i转换为time.Time类型
func Time(i interface{}, format...string) time.Time {
    return GTime(i, format...).Time
}

// 将变量i转换为time.Time类型
func TimeDuration(i interface{}) time.Duration {
    return time.Duration(Int64(i))
}

// 将变量i转换为time.Time类型, 自动识别i为时间戳或者标准化的时间字符串。
func GTime(i interface{}, format...string) *gtime.Time {
    s := String(i)
    if len(s) == 0 {
        return gtime.New()
    }
    // 优先使用用户输入日期格式进行转换
    if len(format) > 0 {
        t, _ := gtime.StrToTimeFormat(s, format[0])
        return t
    }
    if gstr.IsNumeric(s) {
        return gtime.NewFromTimeStamp(Int64(s))
    } else {
        t, _ := gtime.StrToTime(s)
        return t
    }
}
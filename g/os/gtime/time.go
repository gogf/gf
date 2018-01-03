// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 时间管理
package gtime

import (
    "time"
)

// 类似与js中的SetTimeout，一段时间后执行回调函数
func SetTimeout(t time.Duration, callback func()) {
    go func() {
        time.Sleep(t)
        callback()
    }()
}

// 类似与js中的SetInterval，每隔一段时间后执行回调函数，当回调函数返回true，那么继续执行，否则终止执行，该方法是异步的
// 注意：由于采用的是循环而不是递归操作，因此间隔时间将会以上一次回调函数执行完成的时间来计算
func SetInterval(t time.Duration, callback func() bool) {
    go func() {
        for {
            time.Sleep(t)
            if !callback() {
                break
            }
        }
    }()
}

// 获取当前的纳秒数
func Nanosecond() int64 {
    return time.Now().UnixNano()
}

// 获取当前的微秒数
func Microsecond() int64 {
    return time.Now().UnixNano()/1e3
}

// 获取当前的毫秒数
func Millisecond() int64 {
    return time.Now().UnixNano()/1e6
}

// 获取当前的秒数(时间戳)
func Second() int64 {
    return time.Now().UnixNano()/1e9
}

// 获得当前的日期(例如：2006-01-02)
func Date() string {
    return time.Now().Format("2006-01-02")
}

// 获得当前的时间(例如：2006-01-02 15:04:05)
func Datetime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

// 时间戳转换为指定格式的字符串，format格式形如：2006-01-02 03:04:05 PM
// 第二个参数指定需要格式化的时间戳，为非必需参数，默认为当前时间戳
func Format(format string, timestamps...int64) string {
    timestamp := Second()
    if len(timestamps) > 0 {
        timestamp = timestamps[0]
    }
    return time.Unix(timestamp, 0).Format(format)
}

// 字符串转换为时间戳，需要给定字符串时间格式，format格式形如：2006-01-02 03:04:05 PM
func StrToTime(format string, timestr string) (int64, error) {
    t, err := time.Parse(format, timestr)
    if err != nil {
        return 0, err
    }
    return t.Unix(), nil
}

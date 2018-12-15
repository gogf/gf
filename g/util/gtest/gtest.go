// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtest provides useful test utils.
// 测试模块.
package gtest

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/util/gconv"
    "os"
)

// 断言判断
func Assert(value, expect interface{}) {
    if gconv.String(value) != gconv.String(expect) {
        glog.Backtrace(true, 1).Printfln(`[ASSERT] VALUE: %v, EXPECT: %v`, value, expect)
        os.Exit(1)
    }
}

// 提示错误并退出
func Fatal(message...interface{}) {
    glog.Backtrace(true, 1).Println(`[FATAL] `, fmt.Sprint(message...))
    os.Exit(1)
}
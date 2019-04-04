// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/internal/empty"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/util/gutil"
)

// NewVar creates a *Var.
//
// 动态变量
func NewVar(i interface{}, unsafe...bool) *Var {
    return gvar.New(i, unsafe...)
}

// Wait blocks until all the web servers shutdown.
//
// 阻塞等待HTTPServer执行完成(同一进程多HTTPServer情况下)
func Wait() {
    ghttp.Wait()
}

// Dump dumps a variable to stdout with more manually readable.
//
// 格式化打印变量.
func Dump(i...interface{}) {
    gutil.Dump(i...)
}

// Export exports a variable to string with more manually readable.
//
// 格式化导出变量.
func Export(i...interface{}) string {
    return gutil.Export(i...)
}

// Throw throws a exception, which can be caught by Catch function.
// It always be used in TryCatch function.
//
// 抛出一个异常
func Throw(exception interface{}) {
    gutil.Throw(exception)
}

// TryCatch does the try...catch... logic.
func TryCatch(try func(), catch ... func(exception interface{})) {
    gutil.TryCatch(try, catch...)
}

// IsEmpty checks given value empty or not.
// false: integer(0), bool(false), slice/map(len=0), nil;
// true : other.
//
// 判断给定的变量是否为空。
// 整型为0, 布尔为false, slice/map长度为0, 其他为nil的情况，都为空。
// 为空时返回true，否则返回false。
func IsEmpty(value interface{}) bool {
    return empty.IsEmpty(value)
}
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/util/gutil"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/container/gvar"
)

const (
    LOG_LEVEL_ALL  = glog.LEVEL_ALL
    LOG_LEVEL_DEBU = glog.LEVEL_DEBU
    LOG_LEVEL_INFO = glog.LEVEL_INFO
    LOG_LEVEL_NOTI = glog.LEVEL_NOTI
    LOG_LEVEL_WARN = glog.LEVEL_WARN
    LOG_LEVEL_ERRO = glog.LEVEL_ERRO
    LOG_LEVEL_CRIT = glog.LEVEL_CRIT
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
// 打印变量
func Dump(i...interface{}) {
    gutil.Dump(i...)
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
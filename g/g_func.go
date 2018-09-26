// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package g

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/os/glog"
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

// 阻塞等待HTTPServer执行完成(同一进程多HTTPServer情况下)
func Wait() {
    ghttp.Wait()
}

// 是否显示调试信息
func SetDebug(debug bool) {
    glog.SetDebug(debug)
}

// 设置日志的显示等级
func SetLogLevel(level int) {
    glog.SetLevel(level)
}

// 打印变量
func Dump(i...interface{}) {
    gutil.Dump(i...)
}

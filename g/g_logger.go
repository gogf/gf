// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/os/glog"
)

// Disable/Enabled debug of logging globally.
//
// 是否显示调试信息
func SetDebug(debug bool) {
    glog.SetDebug(debug)
}

// Set the logging level globally.
//
// 设置日志的显示等级
func SetLogLevel(level int) {
    glog.SetLevel(level)
}

// Get the global logging level.
//
// 获取设置的日志显示等级
func GetLogLevel() int {
    return glog.GetLevel()
}
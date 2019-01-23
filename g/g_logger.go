// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package g

import (
    "gitee.com/johng/gf/g/os/glog"
)

// 是否显示调试信息
func SetDebug(debug bool) {
    glog.SetDebug(debug)
}

// 设置日志的显示等级
func SetLogLevel(level int) {
    glog.SetLevel(level)
}

// 获取设置的日志显示等级
func GetLogLevel() int {
    return glog.GetLevel()
}
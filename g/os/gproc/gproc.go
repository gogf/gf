// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理.
package gproc

import (
    "os"
    "gitee.com/johng/gf/g/util/gconv"
)

const (
    gPROC_ENV_KEY_PPID_KEY = "gproc.ppid"
)

// 获取当前进程ID
func Pid() int {
    return os.Getpid()
}

// 获取父进程ID
func Ppid() int {
    return gconv.Int(os.Getenv(gPROC_ENV_KEY_PPID_KEY))
}

// 判断当前进程是否为gproc创建的子进程
func IsChild() bool {
    return os.Getenv(gPROC_ENV_KEY_PPID_KEY) != ""
}


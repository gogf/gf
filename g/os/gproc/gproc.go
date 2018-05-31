// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理/通信.
// 本进程管理从syscall, os.StartProcess, exec.Cmd都使用过，
// 最后采用了exec.Cmd来实现多进程管理，这是一个顶层的跨平台封装，兼容性更好，另外两个是偏底层的接口。
package gproc

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/util/gconv"
    "strings"
)

const (
    gPROC_ENV_KEY_PPID_KEY = "gproc.ppid"
    gPROC_TEMP_DIR_ENV_KEY = "gproc.tempdir"
)

// 进程开始执行时间
var processStartTime = time.Now()

// 获取当前进程ID
func Pid() int {
    return os.Getpid()
}

// 获取父进程ID(gproc父进程，如果当前进程本身就是父进程，那么返回自身的pid，不存在时则使用系统父进程)
func PPid() int {
    if !IsChild() {
        return Pid()
    }
    // gPROC_ENV_KEY_PPID_KEY为gproc包自定义的父进程
    ppidValue := os.Getenv(gPROC_ENV_KEY_PPID_KEY)
    if ppidValue != "" {
        return gconv.Int(ppidValue)
    }
    return PPidOS()
}

// 获取父进程ID(系统父进程)
func PPidOS() int {
    return os.Getppid()
}

// 判断当前进程是否为gproc创建的子进程
func IsChild() bool {
    return os.Getenv(gPROC_ENV_KEY_PPID_KEY) != ""
}

// 设置gproc父进程ID，当ppid为0时表示该进程为gproc主进程，否则为gproc子进程
func SetPPid(ppid int) {
    if ppid > 0 {
        os.Setenv(gPROC_ENV_KEY_PPID_KEY, gconv.String(ppid))
    } else {
        os.Unsetenv(gPROC_ENV_KEY_PPID_KEY)
    }
}

// 进程开始执行时间
func StartTime() time.Time {
    return processStartTime
}

// 进程已经运行的时间(毫秒)
func Uptime() int {
    return int(time.Now().UnixNano()/1e6 - processStartTime.UnixNano()/1e6)
}

// 检测环境变量中是否已经存在指定键名
func checkEnvKey(env []string, key string) bool {
    for _, v := range env {
        if len(v) >= len(key) && strings.EqualFold(v[0 : len(key)], key) {
            return true
        }
    }
    return false
}


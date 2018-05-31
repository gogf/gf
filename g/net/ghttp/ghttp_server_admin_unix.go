// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// +build !windows

package ghttp

import (
    "os"
    "syscall"
    "os/signal"
)

// 进程信号量监听消息队列
var procSignalChan = make(chan os.Signal)

// 信号量处理
func handleProcessSignal() {
    var sig os.Signal
    signal.Notify(
        procSignalChan,
        syscall.SIGINT,
        syscall.SIGQUIT,
        syscall.SIGKILL,
        syscall.SIGHUP,
        syscall.SIGTERM,
        syscall.SIGUSR1,
        syscall.SIGUSR2,
    )
    for {
        sig = <- procSignalChan
        switch sig {
        // 进程终止，停止所有子进程运行
        case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM:

            return

            // 用户信号，热重启服务
        case syscall.SIGUSR1:


            // 用户信号，完整重启服务
        case syscall.SIGUSR2:


        default:
        }
    }
}
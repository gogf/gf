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
    "gitee.com/johng/gf/g/os/gproc"
)

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
    )
    for {
        sig = <- procSignalChan
        switch sig {
        // 进程终止，停止所有子进程运行
        case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM:
            sendProcessMsg(gproc.Pid(), gMSG_SHUTDOWN, nil)
            return

            // 用户信号，重启服务
        case syscall.SIGUSR1:
            sendProcessMsg(gproc.Pid(), gMSG_RESTART, nil)

        default:
        }
    }
}
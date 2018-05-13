// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "os"
    "gitee.com/johng/gf/g/os/gproc"
)

// 开启所有Web Server(根据消息启动)
func onCommChildStart(pid int, data []byte) {
    // 进程创建成功之后(开始执行服务时间点为准)，通知主进程自身的存在，并开始执行心跳机制
    sendProcessMsg(gproc.PPid(), gMSG_NEW_FORK, nil)
    // 如果创建自己的父进程非gproc父进程，那么表示该进程为重启创建的进程，创建成功之后需要通知父进程自行销毁
    if gproc.PPidOS() != gproc.PPid() {
        sendProcessMsg(gproc.PPidOS(), gMSG_SHUTDOWN, nil)
        // 在windows下必须等待父进程销毁后才能表明Server资源已被释放，才能开始端口监听，否则会端口资源冲突
        if p, err := os.FindProcess(gproc.PPidOS()); err == nil {
            p.Wait()
        }
    }
    // 开启Web Server服务
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            v.(*Server).startServer(nil)
        }
    })
    // 开始心跳时必须保证主进程时间有值，但是又不能等待主进程消息后再开始检测，因此这里自己更新一下通信时间
    updateProcessChildUpdateTime()
    checkHeartbeat.Set(true)
}
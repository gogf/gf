// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package ghttp

import (
    "gitee.com/johng/gf/g/os/gproc"
)

// 开启所有Web Server(根据消息启动)
func onCommChildStart(pid int, data []byte) {
    // 进程创建成功之后(开始执行服务时间点为准)，通知主进程自身的存在，并开始执行心跳机制
    sendProcessMsg(gproc.PPid(), gMSG_NEW_FORK, nil)
    // 开启Web Server服务
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            v.(*Server).startServer(nil)
        }
    })
    // 开始心跳时必须保证主进程时间有值，但是又不能等待主进程消息后再开始检测，因此这里自己更新一下通信时间
    updateProcessUpdateTime()
    checkHeartbeat.Set(true)
}
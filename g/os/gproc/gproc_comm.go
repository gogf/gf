// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// "不要通过共享内存来通信，而应该通过通信来共享内存"


package gproc

import (
    "os"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gqueue"
)

// 本地进程通信发送消息队列
var commSendQueue    = gqueue.New()
// 本地进程通信接收消息队列
var commReceiveQueue = gqueue.New()
// (用于发送)已建立的PID对应的Conn通信对象
var commPidConnMap   = gmap.NewIntInterfaceMap()

// TCP通信数据结构定义
type Msg struct {
    Pid  int    // PID，来源哪个进程
    Data []byte // 数据
}

// TCP通信数据结构定义
type sendQueueItem struct {
    Pid  int    // PID，发向哪个进程
    Data []byte // 数据
}

// 进程管理/通信初始化操作
func init() {
    go startTcpListening()
}

// 获取指定进程的通信文件地址
func getCommFilePath(pid int) string {
    return getCommDirPath() + gfile.Separator + gconv.String(pid)
}

// 获取进程间通信目录地址
func getCommDirPath() string {
    tempDir := os.Getenv("gproc.tempdir")
    if tempDir == "" {
        tempDir = gfile.TempDir()
    }
    return tempDir + gfile.Separator + "gproc"
}
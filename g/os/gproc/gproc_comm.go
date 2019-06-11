// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// "不要通过共享内存来通信，而应该通过通信来共享内存"


package gproc

import (
    "os"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/container/gmap"
)

// 本地进程通信接收消息队列(按照分组进行构建的map，键值为*gqueue.Queue对象)
var commReceiveQueues = gmap.NewStrAnyMap()

// (用于发送)已建立的PID对应的Conn通信对象，键值为一个Pool，防止并行使用同一个通信对象造成数据重叠
var commPidConnMap    = gmap.NewIntAnyMap()

// TCP通信数据结构定义
type Msg struct {
    Pid   int     // PID，来源哪个进程
    Data  []byte  // 数据
    Group string  // 分组名称
}

// 获取指定进程的通信文件地址
func getCommFilePath(pid int) string {
    return getCommDirPath() + gfile.Separator + gconv.String(pid)
}

// 获取进程间通信目录地址
func getCommDirPath() string {
    tempDir := os.Getenv(gPROC_TEMP_DIR_ENV_KEY)
    if tempDir == "" {
        tempDir = gfile.TempDir()
    }
    return tempDir + gfile.Separator + "gproc"
}
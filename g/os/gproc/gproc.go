// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 进程管理.
package gproc

import (
    "os"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/net/gtcp"
    "net"
)

const (
    gCOMMUNICATION_MAIN_PORT  = 30000
    gCOMMUNICATION_CHILD_PORT = 40000
    gCHILD_PROCESS_ENV_KEY    = "gf.process.manager.child"
    gCHILD_PROCESS_ENV_STRING = gCHILD_PROCESS_ENV_KEY + "=1"
)

// TCP通信数据结构定义
type Msg struct {
    Pid  int    // PID，哪个进程发送的消息
    Data []byte // 参数，消息附带的参数
}

// 获取其他进程传递到当前进程的消息包，阻塞执行
func GetMsg() *Msg {
    if v := msgQueue.PopFront(); v != nil {
        return v.(*Msg)
    }
    return nil
}

// 判断当前进程是否为gproc创建的子进程
func IsChild() bool {
    return os.Getenv(gCHILD_PROCESS_ENV_KEY) != ""
}


// TCP数据通信处理回调函数
// 数据格式：总长度(32bit) | PID(32bit) | 校验(32bit) | 参数(变长)
func tcpServiceHandler(conn net.Conn) {
    buffer := gtcp.Receive(conn, gtcp.Retry{3, 100})
    msgs   := bufferToMsgs(buffer)
    if len(msgs) == 0 {
        conn.Close()
        return
    }
    for _, msg := range msgs {
        msgQueue.PushBack(msg)
    }
}

// 数据解包，防止黏包
func bufferToMsgs(buffer []byte) []*Msg {
    s    := 0
    msgs := make([]*Msg, 0)
    for s < len(buffer) {
        length := gbinary.DecodeToInt(buffer[s : 4])
        if length < 0 || length > len(buffer) {
            s++
            continue
        }
        checksum1 := gbinary.DecodeToUint32(buffer[s + 8 : s + 12])
        checksum2 := gtcp.Checksum(buffer[s + 12 : s + length])
        if checksum1 != checksum2 {
            s++
            continue
        }
        msgs = append(msgs, &Msg {
            Pid  : gbinary.DecodeToInt(buffer[s + 4 : s + 8]),
            Data : buffer[s + 12 : s + length],
        })
        s += length
    }
    return msgs
}

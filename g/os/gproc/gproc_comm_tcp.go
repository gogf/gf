// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// "不要通过共享内存来通信，而应该通过通信来共享内存"


package gproc

import (
    "fmt"
    "net"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

const (
    gPROC_DEFAULT_TCP_PORT = 10000
)

// 创建本地进程TCP通信服务
func startTcpListening() {
    var listen *net.TCPListener
    for i := gPROC_DEFAULT_TCP_PORT; ; i++ {
        addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("127.0.0.1:%d", i))
        if err != nil {
            continue
        }
        listen, err = net.ListenTCP("tcp", addr)
        if err != nil {
            continue
        }
        // 将监听的端口保存到通信文件中(字符串类型存放)
        gfile.PutContents(getCommFilePath(Pid()), gconv.String(i))
        break
    }
    for  {
        if conn, err := listen.Accept(); err != nil {
            glog.Error(err)
        } else if conn != nil {
            go tcpServiceHandler(conn)
        }
    }
}

// TCP数据通信处理回调函数
// 数据格式：总长度(32bit) | PID(32bit) | 校验(32bit) | 参数(变长)
func tcpServiceHandler(conn net.Conn) {
    for {
        if buffer, err := gtcp.Receive(conn, gtcp.Retry{3, 10}); err == nil {
            if len(buffer) > 0 {
                for _, v := range bufferToMsgs(buffer) {
                    commReceiveQueue.PushBack(v)
                }
            }
        } else {
            fmt.Println(err)
            conn.Close()
            return
        }
    }
}

// 数据解包，防止黏包
func bufferToMsgs(buffer []byte) []*Msg {
    s    := 0
    msgs := make([]*Msg, 0)
    for s < len(buffer) {
        length := gbinary.DecodeToInt(buffer[s : s + 4])
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

// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// Web Server进程间通信

package ghttp

import (
    "os"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "fmt"
)

const (
    gMSG_START    = iota
    gMSG_RESTART
    gMSG_SHUTDOWN
    gMSG_EXIT
)


// 处理进程间消息
// 数据格式： 操作(8bit) | 参数(变长)
func (s *Server) handleProcessMsg() {
    for {
        if msg := gproc.Receive(); msg != nil {
            fmt.Println(msg)
            act  := gbinary.DecodeToInt(msg.Data[0 : 1])
            data := msg.Data[1 : ]
            if gproc.IsChild() {
                switch act {
                    case gMSG_START:
                        s.startServer(s.bufferToFdMap(data))
                    case gMSG_RESTART:
                    case gMSG_SHUTDOWN: s.Shutdown()
                    case gMSG_EXIT:     os.Exit(0)

                }
            } else {
                switch act {
                    case gMSG_START:
                    case gMSG_RESTART:
                    case gMSG_SHUTDOWN:
                    case gMSG_EXIT:     os.Exit(0)

                }
            }
        }
    }
}

// 向进程发送操作消息
func (s *Server) sendMsg(pid int, act int, data []byte) {
    gproc.Send(pid, append(gbinary.EncodeInt8(int8(act)), data...))
}


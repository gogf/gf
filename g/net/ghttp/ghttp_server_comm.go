// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// Web Server进程间通信

package ghttp

import (
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "fmt"
    "gitee.com/johng/gf/g/encoding/gjson"
    "os"
)

const (
    gMSG_START    = iota
    gMSG_RESTART
    gMSG_SHUTDOWN
)


// 处理进程间消息
// 数据格式： 操作(8bit) | 参数(变长)
func handleProcessMsg() {
    for {
        if msg := gproc.Receive(); msg != nil {
            fmt.Println(msg)
            act  := gbinary.DecodeToInt(msg.Data[0 : 1])
            data := msg.Data[1 : ]
            if gproc.IsChild() {
                switch act {
                    // 开启所有Web Server(根据消息启动)
                    case gMSG_START:
                        sfm := bufferToServerFdMap(data)
                        for k, v := range sfm {
                            GetServer(k).startServer(v)
                        }

                    // 子进程收到重启消息，那么将自身的ServerFdMap信息收集后发送给主进程，由主进程进行统一调度
                    case gMSG_RESTART:
                        b, _ := gjson.Encode(getServerFdMap())
                        sendProcessMsg(gproc.Ppid(), gMSG_RESTART, b)

                    // 友好关闭服务链接并退出
                    case gMSG_SHUTDOWN:
                        serverMapping.RLockFunc(func(m map[string]interface{}) {
                            for _, v := range m {
                                v.(*Server).Shutdown()
                            }
                        })
                        return

                }
            } else {
                switch act {
                    // 开启服务
                    case gMSG_START:
                        p := procManager.NewProcess(os.Args[0], os.Args, os.Environ())
                        p.Run()
                        sendProcessMsg(p.Pid(), gMSG_START, nil)

                    // 重启服务
                    case gMSG_RESTART:
                        // 创建新的服务进程，使用文件描述来监听同样的端口
                        p := procManager.NewProcess(os.Args[0], os.Args, os.Environ())
                        p.Run()
                        sendProcessMsg(p.Pid(), gMSG_START, data)
                        // 关闭旧的服务进程
                        sendProcessMsg(msg.Pid, gMSG_SHUTDOWN, nil)

                    // 关闭服务
                    case gMSG_SHUTDOWN:
                        procManager.Send(formatMsgBuffer(gMSG_SHUTDOWN, nil))

                }
            }
        }
    }
}

// 向进程发送操作消息
func sendProcessMsg(pid int, act int, data []byte) {
    gproc.Send(pid, formatMsgBuffer(act, data))
}

// 生成一条满足Web Server进程通信协议的消息
func formatMsgBuffer(act int, data []byte) []byte {
    return append(gbinary.EncodeInt8(int8(act)), data...)
}


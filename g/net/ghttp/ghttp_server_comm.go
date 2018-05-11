// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// Web Server进程间通信

package ghttp

import (
    "os"
    "syscall"
    "os/signal"
    "gitee.com/johng/gf/g/os/gproc"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/encoding/gbinary"
)

const (
    gMSG_START       = 10
    gMSG_RESTART     = 20
    gMSG_SHUTDOWN    = 30
    gMSG_NEW_FORK    = 40
    gMSG_REMOVE_PROC = 50
    gMSG_HEARTBEAT   = 60

    gPROC_HEARTBEAT_INTERVAL = 1000 // (毫秒)进程间心跳间隔
    gPROC_HEARTBEAT_TIMEOUT  = 3000 // (毫秒)进程间心跳超时时间，如果子进程在这段内没有接收到任何心跳，那么自动退出，防止可能出现的僵尸子进程
)

// 进程信号量监听消息队列
var procSignalChan    = make(chan os.Signal)

// (主子进程)在第一次创建子进程成功之后才会开始心跳检测，同理对应超时时间才会生效
var heartbeatStarted  = gtype.NewBool()

// 处理进程信号量监控以及进程间消息通信
func handleProcessMsgAndSignal() {
    go handleProcessSignal()
    if gproc.IsChild() {
        go handleChildProcessHeartbeat()
    } else {
        go handleMainProcessHeartbeat()
    }
    handleProcessMsg()
}

// 处理进程间消息
// 数据格式： 操作(8bit) | 参数(变长)
func handleProcessMsg() {
    for {
        if msg := gproc.Receive(); msg != nil {
            // 记录消息日志，用于调试
            //gfile.PutContentsAppend("/tmp/gproc-log",
            //    gconv.String(msg.Pid) + "=>" + gconv.String(gproc.Pid()) + ":" + fmt.Sprintf("%v\n", msg.Data),
            //)
            act  := gbinary.DecodeToUint(msg.Data[0 : 1])
            data := msg.Data[1 : ]
            if gproc.IsChild() {
                // 子进程
                switch act {
                    case gMSG_START:     onCommChildStart(msg.Pid, data)
                    case gMSG_RESTART:   onCommChildRestart(msg.Pid, data)
                    case gMSG_HEARTBEAT: onCommChildHeartbeat(msg.Pid, data)
                    case gMSG_SHUTDOWN:
                        onCommChildShutdown(msg.Pid, data)
                        return
                }
            } else {
                // 父进程
                switch act {
                    case gMSG_START:     onCommMainStart(msg.Pid, data)
                    case gMSG_RESTART:   onCommMainRestart(msg.Pid, data)
                    case gMSG_NEW_FORK:  onCommMainNewFork(msg.Pid, data)
                    case gMSG_HEARTBEAT: onCommMainHeartbeat(msg.Pid, data)
                    case gMSG_REMOVE_PROC:
                        onCommMainRemoveProc(msg.Pid, data)
                        // 如果所有子进程都退出，那么主进程也主动退出
                        if procManager.Size() == 0 {
                            return
                        }
                    case gMSG_SHUTDOWN:
                        onCommMainShutdown(msg.Pid, data)
                        return
                }
            }
        }
    }
}

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

// 向进程发送操作消息
func sendProcessMsg(pid int, act int, data []byte) {
    gproc.Send(pid, formatMsgBuffer(act, data))
}

// 生成一条满足Web Server进程通信协议的消息
func formatMsgBuffer(act int, data []byte) []byte {
    return append(gbinary.EncodeUint8(uint8(act)), data...)
}

// 获取所有Web Server的文件描述符map
func getServerFdMap() map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for k, v := range m {
            sfm[k] = v.(*Server).getListenerFdMap()
        }
    })
    return sfm
}

// 二进制转换为FdMap
func bufferToServerFdMap(buffer []byte) map[string]listenerFdMap {
    sfm := make(map[string]listenerFdMap)
    if len(buffer) > 0 {
        j, _ := gjson.LoadContent(buffer, "json")
        for k, _ := range j.ToMap() {
            m := make(map[string]string)
            for k, v := range j.GetMap(k) {
                m[k] = gconv.String(v)
            }
            sfm[k] = m
        }
    }
    return sfm
}
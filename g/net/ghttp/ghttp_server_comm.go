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
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/encoding/gbinary"
    "gitee.com/johng/gf/g/os/gtime"
)

const (
    gMSG_START       = 10
    gMSG_RELOAD      = 20
    gMSG_RESTART     = 30
    gMSG_SHUTDOWN    = 40
    gMSG_CLOSE       = 45
    gMSG_NEW_FORK    = 50
    gMSG_HEARTBEAT   = 70

    gPROC_HEARTBEAT_INTERVAL    = 1000       // (毫秒)进程间心跳间隔
    gPROC_HEARTBEAT_TIMEOUT     = 3000       // (毫秒)进程间心跳超时时间，如果子进程在这段内没有接收到任何心跳，那么自动退出，防止可能出现的僵尸子进程
)

// 进程信号量监听消息队列
var procSignalChan = make(chan os.Signal)

// 上一次进程间心跳的时间戳
var lastUpdateTime = gtype.NewInt()

// (主子进程)在第一次创建子进程成功之后才会开始心跳检测，同理对应超时时间才会生效
var checkHeartbeat = gtype.NewBool()

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
            //content := gconv.String(msg.Pid) + "=>" + gconv.String(gproc.Pid()) + ":" + fmt.Sprintf("%v\n", msg.Data)
            //fmt.Print(content)
            //gfile.PutContentsAppend("/tmp/gproc-log", content)
            act  := gbinary.DecodeToUint(msg.Data[0 : 1])
            data := msg.Data[1 : ]
            if msg.Pid != gproc.Pid() {
                updateProcessUpdateTime()
            }
            if gproc.IsChild() {
                // ===============
                // 子进程
                // ===============
                switch act {
                    case gMSG_START:     onCommChildStart(msg.Pid, data)
                    case gMSG_RELOAD:    onCommChildReload(msg.Pid, data)
                    case gMSG_RESTART:   onCommChildRestart(msg.Pid, data)
                    case gMSG_CLOSE:     onCommChildClose(msg.Pid, data)
                    case gMSG_HEARTBEAT: onCommChildHeartbeat(msg.Pid, data)
                    case gMSG_SHUTDOWN:  onCommChildShutdown(msg.Pid, data)
                }
            } else {
                // ===============
                // 父进程
                // ===============
                // 任何进程消息都会自动更新最后通信时间记录
                if msg.Pid != gproc.Pid() {
                    updateProcessCommTime(msg.Pid)
                }
                switch act {
                    case gMSG_START:     onCommMainStart(msg.Pid, data)
                    case gMSG_RELOAD:    onCommMainReload(msg.Pid, data)
                    case gMSG_RESTART:   onCommMainRestart(msg.Pid, data)
                    case gMSG_NEW_FORK:  onCommMainNewFork(msg.Pid, data)
                    case gMSG_HEARTBEAT: onCommMainHeartbeat(msg.Pid, data)
                    case gMSG_SHUTDOWN:
                        onCommMainShutdown(msg.Pid, data)
                        return
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

// 关优雅闭进程所有端口的Web Server服务
// 注意，只是关闭Web Server服务，并不是退出进程
func shutdownWebServers() {
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            for _, s := range v.(*Server).servers {
                s.shutdown()
            }
        }
    })
}

// 强制关闭进程所有端口的Web Server服务
// 注意，只是关闭Web Server服务，并不是退出进程
func closeWebServers() {
    serverMapping.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            for _, s := range v.(*Server).servers {
                s.close()
            }
        }
    })
}

// 更新上一次进程间通信的时间
func updateProcessUpdateTime() {
    lastUpdateTime.Set(int(gtime.Millisecond()))
}
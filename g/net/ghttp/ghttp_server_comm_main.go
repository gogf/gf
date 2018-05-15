// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// Web Server进程间通信 - 主进程.
// 管理子进程按照规则听话玩，不然有一百种方法让子进程在本地混不下去.

package ghttp

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gmap"
)

// (主进程)子进程与主进程第一次通信的时间映射map
// 用以识别子进程创建时间先后顺序，当存在多个子进程时，主动销毁旧的子进程
var procFirstTimeMap  = gmap.NewIntIntMap()

// (主进程)主进程与子进程上一次活跃时间映射map
var procUpdateTimeMap = gmap.NewIntIntMap()

// 开启服务
func onCommMainStart(pid int, data []byte) {
    p := procManager.NewProcess(os.Args[0], os.Args, os.Environ())
    p.Start()
    sendProcessMsg(p.Pid(), gMSG_START, nil)
}

// 心跳处理
func onCommMainHeartbeat(pid int, data []byte) {
    updateProcessCommTime(pid)
}

// 重启服务
func onCommMainRestart(pid int, data []byte) {
    // 向所有子进程发送重启命令，子进程将会搜集Web Server信息发送给父进程进行协调重启工作
    procManager.Send(formatMsgBuffer(gMSG_RESTART, nil))
}

// 新建子进程通知
func onCommMainNewFork(pid int, data []byte) {
    procManager.AddProcess(pid)
    checkHeartbeat.Set(true)
}

// 销毁子进程通知
func onCommMainRemoveProc(pid int, data []byte) {
    procManager.RemoveProcess(pid)
}

// 关闭服务，通知所有子进程退出
func onCommMainShutdown(pid int, data []byte) {
    procManager.Send(formatMsgBuffer(gMSG_SHUTDOWN, nil))
}

// 更新指定进程的通信时间记录
func updateProcessCommTime(pid int) {
    procUpdateTimeMap.Set(pid, int(gtime.Millisecond()))
}

// 主进程与子进程相互异步方式发送心跳信息，保持活跃状态
func handleMainProcessHeartbeat() {
    for {
        time.Sleep(gPROC_HEARTBEAT_INTERVAL*time.Millisecond)
        procManager.Send(formatMsgBuffer(gMSG_HEARTBEAT, nil))
        // 清理过期进程
        if checkHeartbeat.Val() {
            for _, pid := range procManager.Pids() {
                if int(gtime.Millisecond()) - procUpdateTimeMap.Get(pid) > gPROC_HEARTBEAT_TIMEOUT {
                    // 这里需要手动从进程管理器中去掉该进程
                    procManager.RemoveProcess(pid)
                    sendProcessMsg(pid, gMSG_SHUTDOWN, nil)
                }
            }
        }
    }
}

// 清理多余的子进程
//func handleMainProcessChildClear() {
//    for {
//        time.Sleep(gPROC_MULTI_CHILD_CLEAR_INTERVAL*time.Millisecond)
//        if procManager.Size() > 1 {
//            minPid  := 0
//            minTime := int(gtime.Millisecond())
//            for _, pid := range procManager.Pids() {
//                if t := procFirstTimeMap.Get(pid); t < minTime {
//                    minPid  = pid
//                    minTime = t
//                }
//            }
//            if minPid > 0 && procUpdateTimeMap.Get(minPid) - procFirstTimeMap.Get(minPid) > gPROC_MULTI_CHILD_CLEAR_MIN_EXPIRE {
//                sendProcessMsg(minPid, gMSG_SHUTDOWN, nil)
//                glog.Printfln("%d: multi child occurred, shutdown %d", gproc.Pid(), minPid)
//            }
//        }
//    }
//}
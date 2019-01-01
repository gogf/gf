// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "gitee.com/johng/gf/g/container/glist"
    "time"
)

// 延迟添加循环任务，delay参数单位为秒
func (w *Wheel) startLoop() {
    go func() {
        for {
            select {
                case <- w.closed:
                    return
                case t := <- w.ticker.C:
                    if w.status.Val() == STATUS_RUNNING {
                        index := w.index.Val()
                        go w.checkEntries(t, w.slots[index])
                        w.index.Set((index + 1) % w.number)
                    }
            }
        }
    }()
}

// 遍历检查可执行循环任务，并异步执行
func (w *Wheel) checkEntries(t time.Time, l *glist.List) {
    for e := l.Front(); e != nil; e = e.Next() {
        entry := e.Value().(*Entry)
        // 是否已命令停止运行
        if entry.status.Val() == STATUS_CLOSED {
            continue
        }
        // 是否满足运行时间间隔
        if !entry.runnableCheck(t) {
            continue
        }
        // 执行异步运行
        go func(e *glist.Element, l *glist.List) {
            defer func() {
                if err := recover(); err != nil {
                    if err == gPANIC_EXIT {
                        entry.status.Set(STATUS_CLOSED)
                    } else {
                        panic(err)
                    }
                }
                if entry.status.Val() != STATUS_CLOSED {
                    entry.status.Set(STATUS_READY)
                } else {
                    // 异步删除，不受锁机制的影响
                    l.Remove(e)
                }
            }()

            entry.Job()
        }(e, l)

    }
}
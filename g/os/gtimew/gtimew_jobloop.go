// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimew

import (
    "container/list"
    "time"
)

// 延迟添加循环任务，delay参数单位为秒
func (w *Wheel) startLoop() {
    go func() {
        for w.status.Val() != STATUS_CLOSED {
            time.Sleep(time.Second)
            if w.status.Val() == STATUS_RUNNING {
                go w.checkEntries(time.Now())
            }
        }
    }()
}

// 遍历检查可执行循环任务，并异步执行
func (w *Wheel) checkEntries(t time.Time) {
    w.entries.RLockFunc(func(l *list.List) {
        for e := l.Front(); e != nil; e = e.Next() {
            entry := e.Value.(*Entry)
            if entry.Meet(t) {
                // 是否已命令停止运行
                if entry.status.Val() == STATUS_CLOSED {
                    continue
                }
                // 判断任务的运行模式
                switch entry.mode.Val() {
                    // 是否只允许单例运行
                    case MODE_SINGLETON:
                        if  entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                            continue
                        }
                    // 只运行一次的任务
                    case MODE_ONCE:
                        if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                            continue
                        }
                }
                // 执行异步运行
                go func(element *list.Element) {
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
                            w.entries.Remove(element)
                        }
                    }()

                    entry.Job()
                }(e)
            }
        }
    })
}
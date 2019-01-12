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

// 开始循环
func (w *wheel) start() {
    go func() {
        ticker := time.NewTicker(time.Duration(w.intervalMs)*time.Millisecond)
        for {
           select {
               case <- w.closed:
                   ticker.Stop()
                   return

               case <- ticker.C:
                    w.proceed()
           }
        }
    }()
}

// 执行时间轮刻度逻辑
func (w *wheel) proceed() {
    n := w.ticks.Add(1)
    l := w.slots[int(n%w.number)]
    if l.Len() > 0 {
        go func(l *glist.List, nowTicks int64) {
            nowMs := time.Now().UnixNano()/1e6
            for i := l.Len(); i > 0; i-- {
                v := l.PopFront()
                if v == nil {
                    break
                }
                entry := v.(*Entry)
                if entry.Status() == STATUS_CLOSED {
                    continue
                }
                // 是否满足运行条件
                if entry.check(nowTicks, nowMs) {
                    // 异步执行运行
                    go func(entry *Entry) {
                        defer func() {
                            if err := recover(); err != nil {
                                if err != gPANIC_EXIT {
                                    panic(err)
                                } else {
                                    entry.Close()
                                }
                            }
                            if entry.Status() == STATUS_RUNNING {
                                entry.SetStatus(STATUS_READY)
                            }
                        }()
                        entry.job()
                    }(entry)
                }
                // 是否继续添运行
                if entry.status.Val() != STATUS_CLOSED {
                    left := entry.interval - (nowTicks - entry.create)
                    if left <= 0 {
                        left           = entry.interval
                        entry.create   = nowTicks
                        entry.createMs = nowMs
                        entry.updateMs = nowMs
                    }
                    w.slots[(nowTicks + left) % w.number].PushBack(entry)
                }
            }
        }(l, n)
    }
}

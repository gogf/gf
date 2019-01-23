// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer

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
               case <- ticker.C:
                   switch w.timer.status.Val() {
                       case STATUS_RUNNING:
                           w.proceed()

                       case STATUS_STOPPED:
                       case STATUS_CLOSED:
                           ticker.Stop()
                           return
                   }

           }
        }
    }()
}

// 执行时间轮刻度逻辑
func (w *wheel) proceed() {
    n      := w.ticks.Add(1)
    l      := w.slots[int(n%w.number)]
    length := l.Len()
    if length > 0 {
        go func(l *glist.List, nowTicks int64) {
            entry := (*Entry)(nil)
            nowMs := time.Now().UnixNano()/1e6
            for i := length; i > 0; i-- {
                if v := l.PopFront(); v == nil {
                    break
                } else {
                    entry = v.(*Entry)
                }
                // 是否满足运行条件
                runnable, addable := entry.check(nowTicks, nowMs)
                if runnable {
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
                // 是否继续添运行, 滚动任务
                if addable {
                    entry.wheel.timer.doAddEntryByParent(entry.rawIntervalMs, entry)
                }
            }
        }(l, n)
    }
}

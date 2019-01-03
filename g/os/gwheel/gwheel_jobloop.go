// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "gitee.com/johng/gf/g/container/glist"
)

// 延迟添加循环任务，delay参数单位为秒
func (w *Wheel) startLoop() {
    go func() {
        for {
            select {
                case <- w.closed:
                    return

                case t := <- w.ticker.C:
                    // 去掉余数，调整为时间轮间隔整数的时间对象
                    n := t.UnixNano()
                    n -= n%w.interval
                    i := w.index.Val()
                    l := w.slots[i]
                    if l.Len() > 0 {
                        go w.checkEntries(n, l)
                    }
                    w.index.Set((i + 1) % w.number)
            }
        }
    }()
}

// 遍历检查可执行循环任务，并异步执行
func (w *Wheel) checkEntries(n int64, l *glist.List) {
    for e := l.Front(); e != nil; e = e.Next() {
        entry := e.Value().(*Entry)
        // 是否已停止运行, 那么移除
        if entry.Status() == STATUS_CLOSED {
            l.Remove(e)
            continue
        }
        // 是否满足运行条件
        if !entry.runnableCheck(n) {
            continue
        }
        // 异步执行运行
        go func(e *glist.Element, l *glist.List) {
            defer func() {
                if err := recover(); err != nil {
                    if err != gPANIC_EXIT {
                        panic(err)
                    } else {
                        entry.Close()
                    }
                }
                if entry.Status() == STATUS_CLOSED {
                    l.Remove(e)
                }
            }()
            entry.Job()
        }(e, l)

    }
}
// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "container/list"
    "gitee.com/johng/gf/g/container/glist"
)

// 开始循环
func (w *wheel) start() {
    go func() {
        for {
           select {
               case <- w.closed:
                   return

               case <- w.ticker.C:
                    w.proceed()
           }
        }
    }()
}

// 执行时间轮刻度逻辑, 遍历检查可执行循环任务，并异步执行
func (w *wheel) proceed() {
    n := w.ticks.Add(1)
    l := w.slots[int(n%int64(w.number))]
    if l.Len() > 0 {
        go func(l *glist.List, ticks int64) {
            l.RLockFunc(func(list *list.List) {
                for e := list.Front(); e != nil; e = e.Next() {
                    entry := e.Value.(*Entry)
                    // 是否满足运行条件
                    if !entry.runnableCheck(ticks) {
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
                            switch entry.Status() {
                                case STATUS_CLOSED:
                                    l.Remove(e)

                                case STATUS_RUNNING:
                                    entry.SetStatus(STATUS_READY)

                            }
                        }()

                        entry.Run()
                    }(e, l)
                }
            })
        }(l, n)
    }
}

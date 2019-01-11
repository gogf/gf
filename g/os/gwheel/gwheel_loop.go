// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "container/list"
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

// 执行时间轮刻度逻辑, 遍历检查可执行循环任务，并异步执行
func (w *wheel) proceed() {
    n := w.ticks.Add(1)
    l := w.slots[int(n%w.number)]
    if l.Len() > 0 {
        go func(l *glist.List, nowTicks int64) {
            nowMs       := time.Now().UnixNano()/1e6
            removeArray := make([]*glist.Element, 0)
            l.RLockFunc(func(list *list.List) {
                for e := list.Front(); e != nil; e = e.Next() {
                    entry := e.Value.(*Entry)
                    // 任务是否已关闭，那么需要删除
                    if entry.Status() == STATUS_CLOSED {
                        removeElement := e
                        removeArray    = append(removeArray, removeElement)
                        continue
                    }
                    // 是否满足运行条件
                    if !entry.runnableCheck(nowTicks, nowMs) {
                        continue
                    }
                    // 异步执行运行
                    go func(entry *Entry, e *glist.Element, l *glist.List) {
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
                                    //fmt.Println("remove:", w.level, entry.interval)
                                    l.Remove(e)

                                case STATUS_RUNNING:
                                    entry.SetStatus(STATUS_READY)
                            }
                        }()

                        //if entry.interval == 8 {
                        //   fmt.Println(w.level, ticks, entry.create, entry.interval,
                        //       entry.times.Val(),
                        //       entry.Status(),
                        //       time.Duration(w.interval),
                        //       time.Duration(entry.interval)*time.Duration(w.interval),
                        //       time.Now(),
                        //       entry.id,
                        //   )
                        //}

                        entry.job()
                    }(entry, e, l)
                }
            })
            if len(removeArray) > 0 {
                l.BatchRemove(removeArray)
            }
        }(l, n)
    }
}

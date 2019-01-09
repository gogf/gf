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


func (w *Wheel) startLoop() {
    go func() {
        for {
           select {
               case <- w.closed:
                   return

               case <- w.ticker.C:
                   n := w.ticks.Add(1)
                   l := w.slots[n%w.number]
                   //if w.interval == 10*time.Millisecond.Nanoseconds() {
                   //   fmt.Println(" loop:", w.ticks.Val(), t, n/1000000)
                   //}
                   if l.Len() > 0 {
                       go w.checkEntries(l, n)
                   }
           }
        }
    }()
}

// 遍历检查可执行循环任务，并异步执行
func (w *Wheel) checkEntries(l *glist.List, ticks int) {
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

}
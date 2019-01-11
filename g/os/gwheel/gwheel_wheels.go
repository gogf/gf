// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 分层时间轮
type Wheels struct {
    levels     []*wheel        // 分层
    length     int             // 层数
    number     int             // 每一层Slot Number
    intervalMs int64           // 最小时间刻度(毫秒)
}

// 创建分层时间轮
func New(slot int, interval time.Duration, level...int) *Wheels {
    length := gDEFAULT_WHEEL_LEVEL
    if len(level) > 0 {
        length = level[0]
    }
    ws := &Wheels {
        levels     : make([]*wheel, length),
        length     : length,
        number     : slot,
        intervalMs : interval.Nanoseconds()/1e6,
    }
    for i := 0; i < length; i++ {
        if i > 0 {
            n           := time.Duration(ws.levels[i - 1].totalMs)*time.Millisecond
            w           := ws.newWheel(i, slot, n)
            ws.levels[i] = w
            ws.levels[i - 1].newEntry(n, w.proceed, false, gDEFAULT_TIMES)
        } else {
            ws.levels[i] = ws.newWheel(i, slot, interval)
        }
    }
    ws.levels[0].start()
    return ws
}

// 创建自定义的循环任务管理对象
func (ws *Wheels) newWheel(level int, slot int, interval time.Duration) *wheel {
    w := &wheel {
        wheels     : ws,
        level      : level,
        slots      : make([]*glist.List, slot),
        number     : int64(slot),
        closed     : make(chan struct{}, 1),
        ticks      : gtype.NewInt64(),
        totalMs    : int64(slot)*interval.Nanoseconds()/1e6,
        createMs   : time.Now().UnixNano()/1e6,
        intervalMs : interval.Nanoseconds()/1e6,
    }
    for i := int64(0); i < w.number; i++ {
        w.slots[i] = glist.New()
    }
    return w
}

// 添加循环任务
func (ws *Wheels) Add(interval time.Duration, job JobFunc) *Entry {
    return ws.newEntry(interval, job, false, gDEFAULT_TIMES)
}

// 添加单例运行循环任务
func (ws *Wheels) AddSingleton(interval time.Duration, job JobFunc) *Entry {
    return ws.newEntry(interval, job, true, gDEFAULT_TIMES)
}

// 添加只运行一次的循环任务
func (ws *Wheels) AddOnce(interval time.Duration, job JobFunc) *Entry {
    return ws.newEntry(interval, job, false, 1)
}

// 添加运行指定次数的循环任务
func (ws *Wheels) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
    return ws.newEntry(interval, job, false, times)
}

// 延迟添加循环任务
func (ws *Wheels) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
    ws.AddOnce(delay, func() {
        ws.Add(interval, job)
    })
}

// 延迟添加单例循环任务
func (ws *Wheels) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
    ws.AddOnce(delay, func() {
        ws.AddSingleton(interval, job)
    })
}

// 延迟添加只运行一次的循环任务
func (ws *Wheels) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
    ws.AddOnce(delay, func() {
        ws.AddOnce(interval, job)
    })
}

// 延迟添加只运行一次的循环任务
func (ws *Wheels) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
    ws.AddOnce(delay, func() {
        ws.AddTimes(interval, times, job)
    })
}

// 关闭分层时间轮
func (ws *Wheels) Close() {
    for _, w := range ws.levels {
        w.Close()
    }
}

// 添加循环任务
func (ws *Wheels) newEntry(interval time.Duration, job JobFunc, singleton bool, times int, from...*wheel) *Entry {
    intervalMs := interval.Nanoseconds()/1e6
    pos, cmp   := ws.binSearchIndex(intervalMs)
    if len(from) > 0 {
        pos = from[0].level - 1
        cmp = -1
    }
    switch cmp {
        // n比最后匹配值小
        case -1:
            i := pos
            for ; i > 0; i-- {
                if intervalMs > ws.levels[i].totalMs {
                    return ws.levels[i].newEntry(interval, job, singleton, times)
                }
            }
            return ws.levels[i].newEntry(interval, job, singleton, times)
        // n比最后匹配值大
        case  1:
            i := pos
            for ; i < ws.length - 1; i++ {
                if intervalMs < ws.levels[i].totalMs {
                    return ws.levels[i].newEntry(interval, job, singleton, times)
                }
            }
            return ws.levels[i].newEntry(interval, job, singleton, times)

        case  0:
            return ws.levels[pos].newEntry(interval, job, singleton, times)
    }
    return nil
}


func (ws *Wheels) binSearchIndex(n int64)(index int, result int) {
    min := 0
    max := ws.length - 1
    mid := 0
    cmp := -2
    for min <= max {
        mid = int((min + max) / 2)
        switch {
            case ws.levels[mid].totalMs == n : cmp =  0
            case ws.levels[mid].totalMs  > n : cmp = -1
            case ws.levels[mid].totalMs  < n : cmp =  1
        }
        switch cmp {
            case -1 : max = mid - 1
            case  1 : min = mid + 1
            case  0 :
                return mid, cmp
        }
    }
    return mid, cmp
}
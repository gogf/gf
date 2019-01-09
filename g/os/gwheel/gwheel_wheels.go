// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "time"
)

// 分层时间轮
type Wheels struct {
    levels    []*wheel        // 分层
    length    int             // 层数
    number    int             // 每一层Slot Number
    interval  int64           // 最小时间刻度(纳秒)
}

// 创建分层时间轮
func New(slot int, interval time.Duration, level...int) *Wheels {
    length := gDEFAULT_WHEEL_LEVEL
    if len(level) > 0 {
        length = level[0]
    }
    w := &Wheels {
        levels   : make([]*wheel, length),
        length   : length,
        number   : slot,
        interval : interval.Nanoseconds(),
    }
    for i := 0; i < length; i++ {
        if i > 0 {
            n          := interval*time.Duration(slot*i)
            w.levels[i] = newWheel(slot, n)
            w.levels[i - 1].Add(n, func() {
                w.levels[i].proceed()
            })
        } else {
            w.levels[i] = newWheel(slot, interval)
        }
    }
    w.levels[0].start()
    return w
}

// 添加循环任务
func (w *Wheels) Add(interval time.Duration, job JobFunc) *Entry {
    n        := interval.Nanoseconds()
    pos, cmp := w.binSearchIndex(n)
    switch cmp {
        case -1 :
            for i := pos; i >= 0; i-- {
                if n > w.levels[i].interval {
                    return w.levels[i].Add(time.Duration(n), job)
                }
            }
            return w.levels[0].Add(time.Duration(n), job)
        case  1 :
            for i := pos; i < w.length; i++ {
                if n > w.levels[i].interval {
                    return w.levels[i].Add(time.Duration(n), job)
                }
            }
            return w.levels[w.length - 1].Add(time.Duration(n), job)
        case  0 :
            return w.levels[pos].Add(time.Duration(n), job)
    }
    return nil
}


func (w *Wheels) binSearchIndex(n int64)(index int, result int) {
    min := 0
    max := w.length - 1
    mid := 0
    cmp := -2
    for min <= max {
        mid = int((min + max) / 2)
        switch {
            case w.levels[mid].interval == n : cmp =  0
            case w.levels[mid].interval  > n : cmp =  1
            case w.levels[mid].interval  < n : cmp = -1
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
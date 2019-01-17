// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 定时器/分层时间轮
type Timer struct {
    status     *gtype.Int      // 状态
    wheels     []*wheel        // 分层
    length     int             // 层数
    number     int             // 每一层Slot Number
    intervalMs int64           // 最小时间刻度(毫秒)
}

// 创建分层时间轮
func New(slot int, interval time.Duration, level...int) *Timer {
    length := gDEFAULT_WHEEL_LEVEL
    if len(level) > 0 {
        length = level[0]
    }
    t := &Timer {
        status     : gtype.NewInt(STATUS_RUNNING),
        wheels     : make([]*wheel, length),
        length     : length,
        number     : slot,
        intervalMs : interval.Nanoseconds()/1e6,
    }
    for i := 0; i < length; i++ {
        if i > 0 {
            n          := time.Duration(t.wheels[i - 1].totalMs)*time.Millisecond
            w          := t.newWheel(i, slot, n)
            t.wheels[i] = w
            t.wheels[i - 1].addEntry(n, w.proceed, false, gDEFAULT_TIMES)
        } else {
            t.wheels[i] = t.newWheel(i, slot, interval)
        }
    }
    t.wheels[0].start()
    return t
}

// 创建自定义的循环任务管理对象
func (t *Timer) newWheel(level int, slot int, interval time.Duration) *wheel {
    w := &wheel {
        timer      : t,
        level      : level,
        slots      : make([]*glist.List, slot),
        number     : int64(slot),
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
func (t *Timer) Add(interval time.Duration, job JobFunc) *Entry {
    return t.addEntry(interval, job, false, gDEFAULT_TIMES)
}

// 添加单例运行循环任务
func (t *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
    return t.addEntry(interval, job, true, gDEFAULT_TIMES)
}

// 添加只运行一次的循环任务
func (t *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
    return t.addEntry(interval, job, false, 1)
}

// 添加运行指定次数的循环任务
func (t *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
    return t.addEntry(interval, job, false, times)
}

// 延迟添加循环任务
func (t *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
    t.AddOnce(delay, func() {
        t.Add(interval, job)
    })
}

// 延迟添加单例循环任务
func (t *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
    t.AddOnce(delay, func() {
        t.AddSingleton(interval, job)
    })
}

// 延迟添加只运行一次的循环任务
func (t *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
    t.AddOnce(delay, func() {
        t.AddOnce(interval, job)
    })
}

// 延迟添加只运行一次的循环任务
func (t *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
    t.AddOnce(delay, func() {
        t.AddTimes(interval, times, job)
    })
}

// 启动定时器
func (t *Timer) Start() {
    t.status.Set(STATUS_RUNNING)
}

// 定制定时器
func (t *Timer) Stop() {
    t.status.Set(STATUS_STOPPED)
}

// 关闭定时器
func (t *Timer) Close() {
    t.status.Set(STATUS_CLOSED)
}

// 添加定时任务
func (t *Timer) addEntry(interval time.Duration, job JobFunc, singleton bool, times int) *Entry {
    intervalMs := interval.Nanoseconds()/1e6
    pos, cmp   := t.binSearchIndex(intervalMs)
    switch cmp {
        // n比最后匹配值小
        case -1:
            i := pos
            for ; i > 0; i-- {
                if intervalMs >= t.wheels[i].intervalMs && intervalMs < t.wheels[i].totalMs {
                    return t.wheels[i].addEntry(interval, job, singleton, times)
                }
            }
            return t.wheels[i].addEntry(interval, job, singleton, times)
        // n比最后匹配值大
        case  1:
            i := pos
            for ; i < t.length - 1; i++ {
                if intervalMs >= t.wheels[i].intervalMs && intervalMs < t.wheels[i].totalMs {
                    return t.wheels[i].addEntry(interval, job, singleton, times)
                }
            }
            return t.wheels[i].addEntry(interval, job, singleton, times)

        case  0:
            return t.wheels[pos].addEntry(interval, job, singleton, times)
    }
    return nil
}

// 二分查找当前任务可以添加的时间轮对象索引
func (t *Timer) binSearchIndex(n int64)(index int, result int) {
    min := 0
    max := t.length - 1
    mid := 0
    cmp := -2
    for min <= max {
        mid = int((min + max) / 2)
        switch {
            case t.wheels[mid].intervalMs == n : cmp =  0
            case t.wheels[mid].intervalMs  > n : cmp = -1
            case t.wheels[mid].intervalMs  < n : cmp =  1
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
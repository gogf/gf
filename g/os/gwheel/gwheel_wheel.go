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

// 循环任务管理对象
type Wheel struct {
    status    *gtype.Int      // 循环任务状态(0: 未执行; 1: 运行中; -1:删除关闭)
    index     *gtype.Int      // 时间轮处理的当前索引位置
    slots     []*glist.List   // 所有的循环任务项, 按照Slot Number进行分组
    number    int             // Slot Number
    closed    chan struct{}   // 停止事件
    create    time.Time       // 创建时间
    ticker    *time.Ticker    // 时间轮间隔
    interval  time.Duration   // 时间间隔(slot时间长度)
}

// 创建使用默认值的时间轮
func NewDefault() *Wheel {
    return New(gDEFAULT_SLOT_NUMBER, gDEFAULT_WHEEL_INTERVAL)
}

// 创建自定义的循环任务管理对象
func New(slot int, interval time.Duration) *Wheel {
    w := &Wheel {
        status    : gtype.NewInt(STATUS_RUNNING),
        index     : gtype.NewInt(),
        slots     : make([]*glist.List, slot),
        number    : slot,
        closed    : make(chan struct{}, 1),
        create    : time.Now(),
        ticker    : time.NewTicker(interval),
        interval  : interval,
    }
    for i := 0; i < w.number; i++ {
        w.slots[i] = glist.New()
    }
    w.startLoop()
    return w
}

// 添加循环任务
func (w *Wheel) Add(interval int, job JobFunc) *Entry {
    return w.newEntry(interval, job, MODE_NORMAL, 0)
}

// 添加单例运行循环任务
func (w *Wheel) AddSingleton(interval int, job JobFunc) *Entry {
    return w.newEntry(interval, job, MODE_SINGLETON, 0)
}

// 添加只运行一次的循环任务
func (w *Wheel) AddOnce(interval int, job JobFunc) *Entry {
    return w.newEntry(interval, job, MODE_ONCE, 0)
}

// 添加运行指定次数的循环任务
func (w *Wheel) AddTimes(interval int, times int, job JobFunc) *Entry {
    return w.newEntry(interval, job, MODE_TIMES, times)
}

// 延迟添加循环任务，delay参数单位为时间轮刻度
func (w *Wheel) DelayAdd(delay int, interval int, job JobFunc) {
    w.AddOnce(delay, func() {
        w.Add(interval, job)
    })
}

// 延迟添加单例循环任务，delay参数单位为时间轮刻度
func (w *Wheel) DelayAddSingleton(delay int, interval int, job JobFunc) {
    w.AddOnce(delay, func() {
        w.AddSingleton(interval, job)
    })
}

// 延迟添加只运行一次的循环任务，delay参数单位为时间轮刻度
func (w *Wheel) DelayAddOnce(delay int, interval int, job JobFunc) {
    w.AddOnce(delay, func() {
        w.AddOnce(interval, job)
    })
}

// 延迟添加只运行一次的循环任务，delay参数单位为时间轮刻度
func (w *Wheel) DelayAddTimes(delay int, interval int, times int, job JobFunc) {
    w.AddOnce(delay, func() {
        w.AddTimes(interval, times, job)
    })
}

// 当前时间轮已注册的任务数
func (w *Wheel) Size() int {
    size := 0
    for _, l := range w.slots {
        size += l.Len()
    }
    return size
}

// 关闭循环任务
func (w *Wheel) Close() {
    w.status.Set(STATUS_CLOSED)
    w.ticker.Stop()
    w.closed <- struct{}{}
}

// 获取所有已注册的循环任务项(按照注册时间从小到大进行排序)
func (w *Wheel) Entries() []*Entry {
    entries := make([]*Entry, 0)
    for _, l := range w.slots {
        for e := l.Front(); e != nil; e = e.Next() {
            entries = append(entries, e.Value().(*Entry))
        }
    }
    return entries
}

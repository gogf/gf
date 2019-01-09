// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务项
type Entry struct {
    singleton *gtype.Bool   // 任务是否单例运行
    status    *gtype.Int    // 任务状态(0: ready;  1: running;  -1: closed)
    times     *gtype.Int64  // 还需运行次数
    create    int64         // 注册时的时间轮ticks
    interval  int64         // 设置的运行间隔(时间轮刻度数量)
    job       JobFunc       // 注册循环任务方法
}

// 任务执行方法
type JobFunc func()

// 创建循环任务
func (w *Wheel) newEntry(interval time.Duration, job JobFunc, singleton bool, times int) *Entry {
    // 安装任务的间隔时间(纳秒)
    n     := interval.Nanoseconds()
    // 计算出所需的插槽数量
    num   := int(n/w.interval)
    if num == 0 {
        // 如果添加的任务间隔时间比时间轮的刻度还小，
        // 那么默认为1个刻度
        num = 1
    }
    ticks := w.ticks.Val()
    entry := &Entry {
        singleton : gtype.NewBool(singleton),
        status    : gtype.NewInt(STATUS_READY),
        times     : gtype.NewInt64(int64(times)),
        job       : job,
        create    : ticks,
        interval  : int64(num),
    }
    // 计算安装的slot数量(可能多个)
    index := int(ticks%int64(w.number))
    for i := 0; i < w.number; i += num {
        w.slots[(i + index + num) % w.number].PushBack(entry)
    }
    return entry
}

// 获取任务状态
func (entry *Entry) Status() int {
    return entry.status.Val()
}

// 设置任务状态
func (entry *Entry) SetStatus(status int) int {
    return entry.status.Set(status)
}

// 关闭当前任务
func (entry *Entry) Close() {
    entry.status.Set(STATUS_CLOSED)
}

// 是否单例运行
func (entry *Entry) IsSingleton() bool {
    return entry.singleton.Val()
}

// 设置单例运行
func (entry *Entry) SetSingleton(enabled bool) {
    entry.singleton.Set(enabled)
}

// 设置任务的运行次数
func (entry *Entry) SetTimes(times int) {
    entry.times.Set(int64(times))
}

// 执行任务
func (entry *Entry) Run() {
    entry.job()
}

// 检测当前任务是否可运行, 参数为当前时间的纳秒数, 精度更高
func (entry *Entry) runnableCheck(ticks int64) bool {
    diff := ticks - entry.create
    if diff > 0 && diff%int64(entry.interval) == 0 {
        // 是否关闭
        if entry.status.Val() == STATUS_CLOSED {
            return false
        }
        // 是否单例
        if entry.IsSingleton() {
            if entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                return false
            }
        }
        // 次数限制
        if entry.times.Add(-1) < 0 {
            if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                return false
            }
        }
        return true
    }
    return false
}

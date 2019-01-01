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
    wheel     *Wheel        // 所属时间轮
    mode      *gtype.Int    // 任务运行模式(0: normal; 1: singleton; 2: once; 3: times)
    status    *gtype.Int    // 循环任务状态(0: ready;  1: running;  -1: stopped)
    times     *gtype.Int    // 还需运行次数, 当mode=3时启用, 当times值为0时表示不再执行(自动该任务删除)
    latest    time.Time     // 任务上一次的运行时间点
    Job       JobFunc       // 注册循环任务方法
    Create    time.Time     // 任务的创建时间点
    Interval  time.Duration // 运行间隔
}

// 任务执行方法
type JobFunc func()

// 创建循环任务
func (w *Wheel) newEntry(interval int, job JobFunc, mode int, times int) *Entry {
    now     := time.Now()
    pos     := (interval + w.index.Val()) % w.number
    entry   := &Entry {
        wheel    : w,
        mode     : gtype.NewInt(mode),
        status   : gtype.NewInt(),
        times    : gtype.NewInt(times),
        latest   : now,
        Job      : job,
        Create   : now,
        Interval : w.interval*time.Duration(interval),
    }
    w.slots[pos].PushBack(entry)
    return entry
}

// 获取任务运行模式
func (entry *Entry) Mode() int {
    return entry.mode.Val()
}

// 设置任务运行模式(0: normal; 1: singleton; 2: once)
func (entry *Entry) SetMode(mode int) {
    entry.mode.Set(mode)
}

// 循环任务状态
func (entry *Entry) Status() int {
    return entry.status.Val()
}

// 启动循环任务
func (entry *Entry) Start() {
    entry.status.Set(STATUS_READY)
}

// 停止循环任务
func (entry *Entry) Stop() {
    entry.status.Set(STATUS_CLOSED)
}

// 检测当前任务是否可运行, 内部将事件转换为微秒数来计算(int64), 精度更高
func (entry *Entry) runnableCheck(t time.Time) bool {
    if t.UnixNano() - entry.latest.UnixNano() >= entry.Interval.Nanoseconds() {
        // 判断任务的运行模式
        switch entry.mode.Val() {
            // 是否只允许单例运行
            case MODE_SINGLETON:
                if  entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                    return false
                }
            // 只运行一次的任务
            case MODE_ONCE:
                if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                    return false
                }
            // 运行指定次数的任务
            case MODE_TIMES:
                if entry.times.Add(-1) <= 0 {
                    if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                        return false
                    }
                }
        }
        entry.latest = t
        return true
    }
    return false
}

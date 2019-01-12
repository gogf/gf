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
    id         int64         // ID
    job        JobFunc       // 注册循环任务方法
    wheel      *wheel        // 所属时间轮
    singleton  *gtype.Bool   // 任务是否单例运行
    status     *gtype.Int    // 任务状态(0: ready;  1: running;  -1: closed)
    times      *gtype.Int    // 还需运行次数
    create     int64         // 注册时的时间轮ticks
    interval   int64         // 设置的运行间隔(时间轮刻度数量)
    createMs   int64         // 创建时间(毫秒)
    updateMs   int64         // 更新时间(上一次检查/运行时间, 毫秒)
    intervalMs int64         // 间隔时间(毫秒)
}

// 任务执行方法
type JobFunc = func()

// 创建循环任务(失败返回nil)
func (w *wheel) addEntry(interval time.Duration, job JobFunc, singleton bool, times int) *Entry {
    ms  := interval.Nanoseconds()/1e6
    num := ms/w.intervalMs
    if num == 0 {
        // 如果安装的任务间隔小于时间轮刻度，
        // 那么将会在下一刻度被执行
        num = 1
    }
    nowMs   := time.Now().UnixNano()/1e6
    ticks   := w.ticks.Val()
    entry   := &Entry {
        id         : time.Now().UnixNano(),
        wheel      : w,
        singleton  : gtype.NewBool(singleton),
        status     : gtype.NewInt(STATUS_READY),
        times      : gtype.NewInt(times),
        job        : job,
        create     : ticks,
        interval   : num,
        createMs   : nowMs,
        updateMs   : nowMs,
        intervalMs : ms,
    }
    // 安装任务
    w.slots[(ticks + num) % w.number].PushBack(entry)
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
    entry.times.Set(times)
}

// 执行任务
func (entry *Entry) Run() {
    entry.job()
}

// 检测当前任务是否可运行, 参数为当前时间的纳秒数, 精度更高
func (entry *Entry) check(nowTicks int64, nowMs int64) bool {
    //if entry.intervalMs == 1400 {
    //   fmt.Println(nowTicks, entry.create, entry.interval, time.Now())
    //}
    // 运行检查
    if diff := nowTicks - entry.create; diff > 0 && diff%entry.interval == 0 {
        // 是否单例
        if entry.IsSingleton() {
            if entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                return false
            }
        }
        // 次数限制
        times := entry.times.Add(-1)
        if times <= 0 {
            entry.status.Set(STATUS_CLOSED)
            if times < 0 {
                return false
            }
        }
        // 是否不限制运行次数
        if times < 2000000000 && times > 1000000000 {
            times = gDEFAULT_TIMES
            entry.times.Set(gDEFAULT_TIMES)
        }
        // 分层转换
        if entry.wheel.level > 0 {
            if diffMs   := nowMs - entry.updateMs; diffMs < entry.intervalMs {
                if leftMs := entry.intervalMs - diffMs; leftMs > entry.wheel.wheels.intervalMs {
                    entry.status.Set(STATUS_CLOSED)
                    delay := time.Duration(leftMs)*time.Millisecond
                    // 往底层添加
                    entry.wheel.wheels.addEntry(delay, entry.job, false, 1)
                    // 延迟重新添加
                    if times > 0 {
                        entry.wheel.wheels.DelayAddTimes(
                            delay,
                            time.Duration(entry.intervalMs)*time.Millisecond,
                            times,
                            entry.job,
                        )
                    }
                    return false
                }
            }
        }
        return true
    }
    return false
}

// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimer

import (
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务项
type Entry struct {
    wheel         *wheel        // 所属时间轮
    job           JobFunc       // 注册循环任务方法
    singleton     *gtype.Bool   // 任务是否单例运行
    status        *gtype.Int    // 任务状态(0: ready;  1: running; 2: stopped; -1: closed), 层级entry共享状态
    times         *gtype.Int    // 还需运行次数
    create        int64         // 注册时的时间轮ticks
    interval      int64         // 设置的运行间隔(时间轮刻度数量)
    createMs      int64         // 创建时间(毫秒)
    intervalMs    int64         // 间隔时间(毫秒)
    rawIntervalMs int64         // 原始间隔
}

// 任务执行方法
type JobFunc = func()

// 创建定时任务。
func (w *wheel) addEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
    ms  := interval.Nanoseconds()/1e6
    num := ms/w.intervalMs
    if num == 0 {
        // 如果安装的任务间隔小于时间轮刻度，
        // 那么将会在下一刻度被执行
        num = 1
    }
    nowMs := time.Now().UnixNano()/1e6
    ticks := w.ticks.Val()
    entry := &Entry {
        wheel         : w,
        job           : job,
        times         : gtype.NewInt(times),
        status        : gtype.NewInt(status),
        create        : ticks,
        interval      : num,
        singleton     : gtype.NewBool(singleton),
        createMs      : nowMs,
        intervalMs    : ms,
        rawIntervalMs : ms,
    }
    // 安装任务
    w.slots[(ticks + num) % w.number].PushBack(entry)
    return entry
}

// 创建定时任务，给定父级Entry, 间隔参数参数为毫秒数.
func (w *wheel) addEntryByParent(interval int64, parent *Entry) *Entry {
    num := interval/w.intervalMs
    if num == 0 {
        num = 1
    }
    nowMs := time.Now().UnixNano()/1e6
    ticks := w.ticks.Val()
    entry := &Entry {
        wheel         : w,
        job           : parent.job,
        times         : parent.times,
        status        : parent.status,
        create        : ticks,
        interval      : num,
        singleton     : parent.singleton,
        createMs      : nowMs,
        intervalMs    : interval,
        rawIntervalMs : parent.rawIntervalMs,
    }
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

// 启动当前任务
func (entry *Entry) Start() {
    entry.status.Set(STATUS_READY)
}

// 停止当前任务
func (entry *Entry) Stop() {
    entry.status.Set(STATUS_STOPPED)
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

// 检测当前任务是否可运行。
func (entry *Entry) check(nowTicks int64, nowMs int64) (runnable, addable bool) {
    switch entry.status.Val() {
        case STATUS_STOPPED:
            return false, true
        case STATUS_CLOSED:
            return false, false
    }
    // 时间轮刻度判断，是否满足运行刻度条件，刻度判断的误差会比较大
    if diff := nowTicks - entry.create; diff > 0 && diff%entry.interval == 0 {
        // 分层转换处理
        if entry.wheel.level > 0 {
            diffMs := nowMs - entry.createMs
            switch {
                // 表示新增(当添加任务后在下一时间轮刻度马上触发)
                case diffMs < entry.wheel.timer.intervalMs:
                    entry.wheel.slots[(nowTicks+entry.interval)%entry.wheel.number].PushBack(entry)
                    return false, false

                // 正常任务
                case diffMs >= entry.wheel.timer.intervalMs:
                    // 任务是否有必要进行分层转换
                    if leftMs := entry.intervalMs - diffMs; leftMs > entry.wheel.timer.intervalMs {

                        // 往底层添加，通过毫秒计算并重新添加任务到对应的时间轮上，减小运行误差
                        entry.wheel.timer.doAddEntryByParent(leftMs, entry)
                        return false, false
                    }
            }
        }
        // 是否单例
        if entry.IsSingleton() {
            // 注意原子操作结果判断
            if entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                return false, true
            }
        }
        // 次数限制
        times := entry.times.Add(-1)
        if times <= 0 {
            // 注意原子操作结果判断
            if entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED || times < 0 {
                return false, false
            }
        }
        // 是否不限制运行次数
        if times < 2000000000 && times > 1000000000 {
            times = gDEFAULT_TIMES
            entry.times.Set(gDEFAULT_TIMES)
        }
        return true, true
    }
    return false, true
}

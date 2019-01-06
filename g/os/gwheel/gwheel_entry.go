// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gwheel

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务项
type Entry struct {
    singleton *gtype.Bool   // 任务是否单例运行
    status    *gtype.Int    // 任务状态(0: ready;  1: running;  -1: closed)
    times     *gtype.Int    // 还需运行次数(<0: 无限制; >=0: 限制次数)
    update    *gtype.Int64  // 任务上一次的运行时间点(纳秒时间戳)
    interval  int64         // 设置的运行间隔(纳秒)
    create    int64         // 创建的时间点(纳秒, 时间轮刻度整数)
    Job       JobFunc       // 注册循环任务方法
    Create    time.Time     // 任务的创建时间点
}

// 任务执行方法
type JobFunc func()

// 创建循环任务
func (w *Wheel) newEntry(interval time.Duration, job JobFunc, singleton bool, times int) (*Entry, error) {
    // 安装任务的间隔时间(纳秒)
    n     := interval.Nanoseconds()
    // 计算出所需的插槽数量
    num   := int(n/w.interval)
    if num == 0 {
      return nil, errors.New(fmt.Sprintf(`interval "%v" should not be less than timing wheel interval "%v"`, interval, time.Duration(w.interval)))
    }
    now   := time.Now().UnixNano()
    entry := &Entry {
        singleton : gtype.NewBool(singleton),
        status    : gtype.NewInt(STATUS_READY),
        times     : gtype.NewInt(times),
        update    : gtype.NewInt64(now - (now%w.interval)),
        Job       : job,
        interval  : n,
    }
    // 计算安装的slot数量(可能多个)
    index := w.index.Val()
    for i := 0; i < w.number; i += num {
        w.slots[(i + index + num) % w.number].PushBack(entry)
    }
    return entry, nil
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

// 检测当前任务是否可运行, 参数为当前时间的纳秒数, 精度更高
func (entry *Entry) runnableCheck(n int64) bool {
    if n - entry.update.Val() >= entry.interval {
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
        if entry.times.Add(-1) == 0 {
            if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                return false
            }
        }
        entry.update.Set(n)
        return true
    }
    return false
}

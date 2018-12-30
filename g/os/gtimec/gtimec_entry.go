// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimec

import (
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务项
type Entry struct {
    mode      *gtype.Int    // 任务运行模式(0: normal; 1: singleton; 2: once)
    status    *gtype.Int    // 循环任务状态(0: ready;  1: running;  -1: stopped)
    Job       func()        // 注册循环任务方法
    Create    int64         // 创建时间戳(秒)
    Interval  int64         // 运行间隔(秒)
}

// 创建循环任务
func newEntry(interval int, job func(), mode int) *Entry {
    return &Entry {
        mode      : gtype.NewInt(mode),
        status    : gtype.NewInt(),
        Job       : job,
        Create    : time.Now().Unix(),
        Interval  : int64(interval),
    }
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

// 给定时间是否满足当前循环任务运行间隔
func (entry *Entry) meet(t time.Time) bool {
    diff := t.Unix() - entry.Create
    if diff > 0 {
        return diff%entry.Interval == 0
    }
    return false
}


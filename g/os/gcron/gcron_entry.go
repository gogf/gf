// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron

import (
    "gitee.com/johng/gf/g/container/gtype"
    "reflect"
    "runtime"
    "time"
)

// 定时任务项
type Entry struct {
    mode      *gtype.Int    // 任务运行模式(0: normal; 1: singleton; 2: once)
    status    *gtype.Int    // 定时任务状态(0: ready;  1: running;  -1: stopped)
    schedule  *cronSchedule // 定时任务配置对象
    Name      string        // 定时任务名称
    Job       func()        // 注册定时任务方法
    JobName   string        // 注册定时任务名称
    Time      time.Time     // 注册时间
}

// 创建定时任务
func newEntry(pattern string, job func(), name ... string) (*Entry, error) {
    schedule, err := newSchedule(pattern)
    if err != nil {
        return nil, err
    }
    entry := &Entry {
        mode      : gtype.NewInt(),
        status    : gtype.NewInt(),
        schedule  : schedule,
        Job       : job,
        JobName   : runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
        Time      : time.Now(),
    }
    if len(name) > 0 {
        entry.Name = name[0]
    }
    return entry, nil
}

// 设置任务运行模式(0: normal; 1: singleton; 2: once)
func (entry *Entry) SetMode(mode int) {
    entry.mode.Set(mode)
}

// 定时任务状态
func (entry *Entry) Status() int {
    return entry.status.Val()
}

// 启动定时任务
func (entry *Entry) Start() {
    entry.status.Set(STATUS_READY)
}

// 停止定时任务
func (entry *Entry) Stop() {
    entry.status.Set(STATUS_CLOSED)
}

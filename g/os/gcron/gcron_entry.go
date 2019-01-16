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
    singleton  *gtype.Bool   // 任务是否单例运行
    times      *gtype.Int    // 还需运行次数
    status     *gtype.Int    // 定时任务状态(0: ready;  1: running;  -1: stopped)
    schedule   *cronSchedule // 定时任务配置对象
    Name       string        // 定时任务名称
    Job        func()        // 注册定时任务方法
    JobName    string        // 注册定时任务名称
    Time       time.Time     // 注册时间
}

// 创建定时任务
func newEntry(pattern string, job func(), singleton bool, times int, name ... string) (*Entry, error) {
    schedule, err := newSchedule(pattern)
    if err != nil {
        return nil, err
    }
    entry := &Entry {
        singleton : gtype.NewBool(singleton),
        times     : gtype.NewInt(times),
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

// 定时任务状态
func (entry *Entry) Status() int {
    return entry.status.Val()
}

// 设置定时任务状态, 返回设置之前的状态
func (entry *Entry) SetStatus(status int) int {
    return entry.status.Set(status)
}

// 启动定时任务
func (entry *Entry) Start() {
    entry.status.Set(STATUS_READY)
}

// 停止定时任务
func (entry *Entry) Stop() {
    entry.status.Set(STATUS_STOPPED)
}

// 关闭定时任务
func (entry *Entry) Close() {
    entry.status.Set(STATUS_CLOSED)
}

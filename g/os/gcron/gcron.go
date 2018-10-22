// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 定时任务.
package gcron

import (
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/third/github.com/robfig/cron"
    "reflect"
    "runtime"
)

// 定时任务项
type Entry struct {
    Spec string      // 注册定时任务时间格式
    Cmd  string      // 注册定时任务名称
    Time *gtime.Time // 注册时间
}

var (
    // 默认的cron管理对象
    defaultCron = cron.New()
    // 当前cron的运行状态(0: 未执行; >0: 运行中)
    cronStatus  = gtype.NewInt()
    // 注册定时任务项
    cronEntries = garray.New(0, 0, true)
)

// 添加执行方法
func Add(spec string, f func()) error {
    // 底层的AddFunc是并发安全的
    err := defaultCron.AddFunc(spec, f)
    if err == nil {
        if cronStatus.Add(1) == 1 {
            go defaultCron.Run()
        }
        cronEntries.Append(Entry{
            Spec : spec,
            Cmd  : runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
            Time : gtime.Now(),
        })
    }
    return err
}

// 获取所有已注册的定时任务项
func Entries() []Entry {
    length  := cronEntries.Len()
    entries := make([]Entry, length)
    for i := 0; i < length; i++ {
        entries[i] = cronEntries.Get(i).(Entry)
    }
    return entries
}

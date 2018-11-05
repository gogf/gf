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
)

// 定时任务项
type Entry struct {
    Spec   string      // 注册定时任务时间格式
    Cmd    string      // 注册定时任务名称
    Time   *gtime.Time // 注册时间
    Name   string      // 定时任务名称
    Status *gtype.Int  // 定时任务状态(0: 未执行; > 0: 运行中)
    cron   *cron.Cron  // 定时任务单独的底层定时管理对象
}

// 定时任务管理对象
type Cron struct {
    cron    *cron.Cron    // 底层定时管理对象
    entries *garray.Array // 定时任务注册项
    status  *gtype.Int    // 默认定时任务管理对象状态(不带名称的定时任务，0: 未执行; > 0: 运行中)
}

var (
    // 默认的cron管理对象
    defaultCron = New()
)

// 创建自定义的定时任务管理对象
func New() *Cron {
    return &Cron {
        cron    : cron.New(),
        entries : garray.New(0, 0, true),
        status  : gtype.NewInt(),
    }
}

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(spec string, f func(), name ... string) error {
    return defaultCron.Add(spec, f, name...)
}

// 延迟添加定时任务，delay参数单位为秒
func DelayAdd(delay int, spec string, f func(), name ... string) {
    defaultCron.DelayAdd(delay, spec, f, name...)
}

// 检索指定名称的定时任务
func Search(name string) *Entry {
    return defaultCron.Search(name)
}

// 根据指定名称删除定时任务
func Remove(name string) {
    defaultCron.Remove(name)
}

// 获取所有已注册的定时任务项
func Entries() []*Entry {
    return defaultCron.Entries()
}

// 启动指定的定时任务
func Start(name string) {
    defaultCron.Start(name)
}

// 停止指定的定时任务
func Stop(name string) {
    defaultCron.Stop(name)
}

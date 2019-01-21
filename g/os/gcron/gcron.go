// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gcron implements a cron pattern parser and job runner.
// 
// 定时任务.
package gcron

import (
    "gitee.com/johng/gf/g/os/gtimer"
    "math"
    "time"
)

const (
    STATUS_READY   = gtimer.STATUS_READY
    STATUS_RUNNING = gtimer.STATUS_RUNNING
    STATUS_STOPPED = gtimer.STATUS_STOPPED
    STATUS_CLOSED  = gtimer.STATUS_CLOSED

    gDEFAULT_TIMES = math.MaxInt32
)

var (
    // 默认的cron管理对象
    defaultCron = New()
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(pattern string, job func(), name ... string) (*Entry, error) {
    return defaultCron.Add(pattern, job, name...)
}

// 添加单例运行定时任务
func AddSingleton(pattern string, job func(), name ... string) (*Entry, error) {
    return defaultCron.AddSingleton(pattern, job, name...)
}

// 添加只运行一次的定时任务
func AddOnce(pattern string, job func(), name ... string) (*Entry, error) {
    return defaultCron.AddOnce(pattern, job, name...)
}

// 添加运行指定次数的定时任务
func AddTimes(pattern string, times int, job func(), name ... string) (*Entry, error) {
    return defaultCron.AddTimes(pattern, times, job, name...)
}

// 延迟添加定时任务
func DelayAdd(delay time.Duration, pattern string, job func(), name ... string) {
    defaultCron.DelayAdd(delay, pattern, job, name...)
}

// 延迟添加单例定时任务，delay参数单位为秒
func DelayAddSingleton(delay time.Duration, pattern string, job func(), name ... string) {
    defaultCron.DelayAddSingleton(delay, pattern, job, name...)
}

// 延迟添加只运行一次的定时任务，delay参数单位为秒
func DelayAddOnce(delay time.Duration, pattern string, job func(), name ... string) {
    defaultCron.DelayAddOnce(delay, pattern, job, name...)
}

// 延迟添加运行指定次数的定时任务，delay参数单位为秒
func DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ... string) {
    defaultCron.DelayAddTimes(delay, pattern, times, job, name...)
}

// 检索指定名称的定时任务
func Search(name string) *Entry {
    return defaultCron.Search(name)
}

// 根据指定名称删除定时任务
func Remove(name string) {
    defaultCron.Remove(name)
}

// 获取所有已注册的定时任务数量
func Size() int {
    return defaultCron.Size()
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

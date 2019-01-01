// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gwheel provides Timing Wheel for interval jobs running and management/时间轮.
// 高效的时间轮任务执行管理，用于管理异步的间隔运行任务，或者异步只运行一次的任务(默认最小时间粒度为秒)。
// 与其他定时任务管理模块的区别是，时间轮模块只管理间隔执行任务，并且更注重执行效率(纳秒级别)。
package gwheel

import "time"

const (
    MODE_NORMAL    = 0
    MODE_SINGLETON = 1
    MODE_ONCE      = 2
    MODE_TIMES     = 3

    STATUS_READY   = 0
    STATUS_RUNNING = 1
    STATUS_CLOSED  = -1

    gPANIC_EXIT    = "exit"

    gDEFAULT_SLOT_NUMBER    = 10
    gDEFAULT_WHEEL_INTERVAL = 100*time.Millisecond
)

var (
    // 默认的wheel管理对象
    defaultWheel = NewDefault()
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(interval int, job JobFunc) *Entry {
    return defaultWheel.Add(10*interval, job)
}

// 添加单例运行循环任务
func AddSingleton(interval int, job JobFunc) *Entry {
    return nil
    return defaultWheel.AddSingleton(10*interval, job)
}

// 添加只运行一次的循环任务
func AddOnce(interval int, job JobFunc) *Entry {
    return defaultWheel.AddOnce(10*interval, job)
}

// 添加运行指定次数的循环任务
func AddTimes(interval int, times int, job JobFunc) *Entry {
    return defaultWheel.AddTimes(10*interval, times, job)
}

// 延迟添加循环任务，delay参数单位为秒
func DelayAdd(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAdd(delay, 10*interval, job)
}

// 延迟添加单例循环任务，delay参数单位为秒
func DelayAddSingleton(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAddSingleton(delay, 10*interval, job)
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func DelayAddOnce(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAddOnce(delay, 10*interval, job)
}

// 延迟添加运行指定次数的循环任务，delay参数单位为秒
func DelayAddTimes(delay int, interval int, times int, job JobFunc) {
    defaultWheel.DelayAddTimes(delay, 10*interval, times, job)
}

// 获取所有已注册的循环任务项
func Entries() []*Entry {
   return defaultWheel.Entries()
}

// 当前时间轮已注册的任务数
func Size() int {
    return defaultWheel.Size()
}

// 在Job方法中调用，停止当前运行的Job
func ExitJob() {
    panic(gPANIC_EXIT)
}

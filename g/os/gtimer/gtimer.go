// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtimer implements Hierarchical Timing Wheel for interval/delayed jobs running and management.
// 
// 任务定时器,
// 高性能的分层时间轮任务管理模块，用于管理间隔/延迟运行任务。
// 与gcron模块的区别是，时间轮模块只管理间隔执行任务，并且更注重执行效率(纳秒级别)。
// 需要注意执行时间间隔的准确性问题: https://github.com/golang/go/issues/14410
package gtimer

import (
    "gitee.com/johng/gf/g/internal/cmdenv"
    "math"
    "time"
)

const (
    STATUS_READY            = 0
    STATUS_RUNNING          = 1
    STATUS_STOPPED          = 2
    STATUS_CLOSED           = -1
    gPANIC_EXIT             = "exit"
    gDEFAULT_TIMES          = math.MaxInt32
    gDEFAULT_SLOT_NUMBER    = 10
    gDEFAULT_WHEEL_INTERVAL = 50
    gDEFAULT_WHEEL_LEVEL    = 6
)

var (
    // 默认定时器属性参数值
    defaultSlots    = cmdenv.Get("gf.gtimer.slots",    gDEFAULT_SLOT_NUMBER).Int()
    defaultLevel    = cmdenv.Get("gf.gtimer.level",    gDEFAULT_WHEEL_LEVEL).Int()
    defaultInterval = cmdenv.Get("gf.gtimer.interval", gDEFAULT_WHEEL_INTERVAL).TimeDuration()*time.Millisecond
    // 默认的wheel管理对象
    defaultTimer    = New(defaultSlots, defaultInterval, defaultLevel)
)

// 类似与js中的SetTimeout，一段时间后执行回调函数。
func SetTimeout(delay time.Duration, job JobFunc) {
    AddOnce(delay, job)
}

// 类似与js中的SetInterval，每隔一段时间执行指定回调函数。
func SetInterval(interval time.Duration, job JobFunc) {
    Add(interval, job)
}

// 添加执行方法。
func Add(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.Add(interval, job)
}

// 添加执行方法，更多参数控制。
func AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
    return defaultTimer.AddEntry(interval, job, singleton, times, status)
}

// 添加单例运行循环任务。
func AddSingleton(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.AddSingleton(interval, job)
}

// 添加只运行一次的循环任务。
func AddOnce(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.AddOnce(interval, job)
}

// 添加运行指定次数的循环任务。
func AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
    return defaultTimer.AddTimes(interval, times, job)
}

// 延迟添加循环任务。
func DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
    defaultTimer.DelayAdd(delay, interval, job)
}

// 延迟添加循环任务, 支持完整的参数。
func DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
    defaultTimer.DelayAddEntry(delay, interval, job, singleton, times, status)
}

// 延迟添加单例循环任务，delay参数单位为秒
func DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
    defaultTimer.DelayAddSingleton(delay, interval, job)
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
    defaultTimer.DelayAddOnce(delay, interval, job)
}

// 延迟添加运行指定次数的循环任务，delay参数单位为秒
func DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
    defaultTimer.DelayAddTimes(delay, interval, times, job)
}

// 在Job方法中调用，停止当前运行的任务。
func Exit() {
    panic(gPANIC_EXIT)
}

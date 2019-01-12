// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtimer implements Levelled Timing Wheel for interval/delayed jobs running and management/任务定时器(分层时间轮).
// 高效的时间轮任务管理模块，用于管理间隔/延迟运行任务。
// 与gcron模块的区别是，时间轮模块只管理间隔执行任务，并且更注重执行效率(纳秒级别)。
// 需要注意执行时间间隔的准确性问题: https://github.com/golang/go/issues/14410
package gtimer

import (
    "math"
    "time"
)

const (
    STATUS_READY            = 0
    STATUS_RUNNING          = 1
    STATUS_CLOSED           = -1
    gPANIC_EXIT             = "exit"
    gDEFAULT_TIMES          = math.MaxInt32
    gDEFAULT_SLOT_NUMBER    = 10
    gDEFAULT_WHEEL_INTERVAL = 50*time.Millisecond
    gDEFAULT_WHEEL_LEVEL    = 6
)

var (
    // 默认的wheel管理对象
    defaultTimer = New(gDEFAULT_SLOT_NUMBER, gDEFAULT_WHEEL_INTERVAL, gDEFAULT_WHEEL_LEVEL)
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.Add(interval, job)
}

// 添加单例运行循环任务
func AddSingleton(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.AddSingleton(interval, job)
}

// 添加只运行一次的循环任务
func AddOnce(interval time.Duration, job JobFunc) *Entry {
    return defaultTimer.AddOnce(interval, job)
}

// 添加运行指定次数的循环任务
func AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
    return defaultTimer.AddTimes(interval, times, job)
}

// 延迟添加循环任务，delay参数单位为秒
func DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
    defaultTimer.DelayAdd(delay, interval, job)
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

// 在Job方法中调用，停止当前运行的任务
func Exit() {
    panic(gPANIC_EXIT)
}

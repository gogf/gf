// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtimew provides Time Wheel for interval jobs running management/时间轮.
// 高效的时间轮任务执行管理，用于管理异步的间隔运行任务，或者异步只运行一次的任务(最小时间粒度为秒)。
// 与其他定时任务管理模块的区别是，时间轮模块只管理间隔执行任务，并且更注重执行效率(纳秒级别)。
package gtimew

const (
    MODE_NORMAL    = 0
    MODE_SINGLETON = 1
    MODE_ONCE      = 2

    STATUS_READY   = 0
    STATUS_RUNNING = 1
    STATUS_CLOSED  = -1

    gPANIC_EXIT    = "exit"
)

var (
    // 默认的wheel管理对象
    defaultWheel = New()
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(interval int, job JobFunc) *Entry {
    return defaultWheel.Add(interval, job)
}

// 添加单例运行循环任务
func AddSingleton(interval int, job JobFunc) *Entry {
    return defaultWheel.AddSingleton(interval, job)
}

// 添加只运行一次的循环任务
func AddOnce(interval int, job JobFunc) *Entry {
    return defaultWheel.AddOnce(interval, job)
}

// 延迟添加循环任务，delay参数单位为秒
func DelayAdd(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAdd(delay, interval, job)
}

// 延迟添加单例循环任务，delay参数单位为秒
func DelayAddSingleton(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAddSingleton(delay, interval, job)
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func DelayAddOnce(delay int, interval int, job JobFunc) {
    defaultWheel.DelayAddOnce(delay, interval, job)
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

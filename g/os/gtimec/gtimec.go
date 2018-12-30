// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gtimec provides Time Circle for interval job running management/时间轮.
// 高效的时间轮任务执行管理，用于管理异步的间隔运行任务，或者异步延迟只运行一次的任务(最小时间粒度为秒)。
// 与其他定时任务管理模块的区别是，时间轮模块只管理间隔执行任务，并且更注重执行效率(纳秒级别)。
package gtimec

const (
    MODE_NORMAL    = 0
    MODE_SINGLETON = 1
    MODE_ONCE      = 2

    STATUS_READY   = 0
    STATUS_RUNNING = 1
    STATUS_CLOSED  = -1
)

var (
    // 默认的circle管理对象
    defaultCircle = New()
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(interval int, job func()) *Entry {
    return defaultCircle.Add(interval, job)
}

// 添加单例运行循环任务
func AddSingleton(interval int, job func()) *Entry {
    return defaultCircle.AddSingleton(interval, job)
}

// 添加只运行一次的循环任务
func AddOnce(interval int, job func()) *Entry {
    return defaultCircle.AddOnce(interval, job)
}

// 延迟添加循环任务，delay参数单位为秒
func DelayAdd(delay int, interval int, job func()) {
    defaultCircle.DelayAdd(delay, interval, job)
}

// 延迟添加单例循环任务，delay参数单位为秒
func DelayAddSingleton(delay int, interval int, job func()) {
    defaultCircle.DelayAddSingleton(delay, interval, job)
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func DelayAddOnce(delay int, interval int, job func()) {
    defaultCircle.DelayAddOnce(delay, interval, job)
}

//// 获取所有已注册的循环任务项
//func Entries() []*Entry {
//    return defaultCircle.Entries()
//}

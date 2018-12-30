// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimec

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务管理对象
type Circle struct {
    status   *gtype.Int  // 循环任务状态(0: 未执行; 1: 运行中; -1:删除关闭)
    entries  *glist.List // 所有的循环任务项
}

// 创建自定义的循环任务管理对象
func New() *Circle {
    circle := &Circle {
        status   : gtype.NewInt(STATUS_RUNNING),
        entries  : glist.New(),
    }
    circle.startLoop()
    return circle
}

// 添加循环任务
func (c *Circle) Add(interval int, job func()) *Entry {
    entry := newEntry(interval, job, MODE_NORMAL)
    c.entries.PushBack(entry)
    return entry
}

// 添加单例运行循环任务
func (c *Circle) AddSingleton(interval int, job func()) *Entry {
    entry := newEntry(interval, job, MODE_SINGLETON)
    c.entries.PushBack(entry)
    return entry
}

// 添加只运行一次的循环任务
func (c *Circle) AddOnce(interval int, job func()) *Entry {
    entry := newEntry(interval, job, MODE_ONCE)
    c.entries.PushBack(entry)
    return entry
}

// 延迟添加循环任务，delay参数单位为秒
func (c *Circle) DelayAdd(delay int, interval int, job func()) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        c.Add(interval, job)
    }()
}

// 延迟添加单例循环任务，delay参数单位为秒
func (c *Circle) DelayAddSingleton(delay int, interval int, job func()) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        c.AddSingleton(interval, job)
    }()
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func (c *Circle) DelayAddOnce(delay int, interval int, job func()) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        c.AddOnce(interval, job)
    }()
}

// 关闭循环任务
func (c *Circle) Close() {
    c.status.Set(STATUS_CLOSED)
}

//// 获取所有已注册的循环任务项(按照注册时间从小到大进行排序)
//func (c *Circle) Entries() []*Entry {
//    array := garray.NewSortedArray(c.entries.Len(), func(v1, v2 interface{}) int {
//        entry1 := v1.(*Entry)
//        entry2 := v2.(*Entry)
//        if entry1.Create > entry2.Create {
//            return 1
//        }
//        return -1
//    }, false)
//    c.entries.RLockFunc(func(m map[string]interface{}) {
//        for _, v := range m {
//            array.Add(v.(*Entry))
//        }
//    })
//    entries := make([]*Entry, array.Len())
//    array.RLockFunc(func(array []interface{}) {
//        for k, v := range array {
//            entries[k] = v.(*Entry)
//        }
//    })
//    return entries
//}

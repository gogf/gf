// Copyright 2019 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtimew

import (
    "container/list"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "time"
)

// 循环任务管理对象
type Wheel struct {
    status   *gtype.Int  // 循环任务状态(0: 未执行; 1: 运行中; -1:删除关闭)
    entries  *glist.List // 所有的循环任务项
}

// 创建自定义的循环任务管理对象
func New() *Wheel {
    wheel := &Wheel {
        status   : gtype.NewInt(STATUS_RUNNING),
        entries  : glist.New(),
    }
    wheel.startLoop()
    return wheel
}

// 添加循环任务
func (w *Wheel) Add(interval int, job JobFunc) *Entry {
    entry := newEntry(interval, job, MODE_NORMAL)
    w.entries.PushBack(entry)
    return entry
}

// 添加单例运行循环任务
func (w *Wheel) AddSingleton(interval int, job JobFunc) *Entry {
    entry := newEntry(interval, job, MODE_SINGLETON)
    w.entries.PushBack(entry)
    return entry
}

// 添加只运行一次的循环任务
func (w *Wheel) AddOnce(interval int, job JobFunc) *Entry {
    entry := newEntry(interval, job, MODE_ONCE)
    w.entries.PushBack(entry)
    return entry
}

// 延迟添加循环任务，delay参数单位为秒
func (w *Wheel) DelayAdd(delay int, interval int, job JobFunc) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        w.Add(interval, job)
    }()
}

// 延迟添加单例循环任务，delay参数单位为秒
func (w *Wheel) DelayAddSingleton(delay int, interval int, job JobFunc) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        w.AddSingleton(interval, job)
    }()
}

// 延迟添加只运行一次的循环任务，delay参数单位为秒
func (w *Wheel) DelayAddOnce(delay int, interval int, job JobFunc) {
    go func() {
        time.Sleep(time.Duration(delay)*time.Second)
        w.AddOnce(interval, job)
    }()
}

// 关闭循环任务
func (w *Wheel) Close() {
    w.status.Set(STATUS_CLOSED)
}

// 获取所有已注册的循环任务项(按照注册时间从小到大进行排序)
func (w *Wheel) Entries() []*Entry {
   array := garray.NewSortedArray(w.entries.Len(), func(v1, v2 interface{}) int {
       entry1 := v1.(*Entry)
       entry2 := v2.(*Entry)
       if entry1.Create > entry2.Create {
           return 1
       }
       return -1
   }, false)
    w.entries.RLockFunc(func(l *list.List) {
        for e := l.Front(); e != nil; e = e.Next() {
            array.Add(e.Value.(*Entry))
        }
   })
   entries := make([]*Entry, array.Len())
   array.RLockFunc(func(array []interface{}) {
       for k, v := range array {
           entries[k] = v.(*Entry)
       }
   })
   return entries
}

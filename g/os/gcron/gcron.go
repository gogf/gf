// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 定时任务.
package gcron

import (
    "errors"
    "fmt"
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
    Name string      // 定时任务名称
    cron *cron.Cron  // 底层定时管理对象
}

var (
    // 默认的cron管理对象
    defaultCron = cron.New()
    // 当前cron的运行状态(0: 未执行; > 0: 运行中)
    cronStatus  = gtype.NewInt()
    // 注册定时任务项
    cronEntries = garray.New(0, 0, true)
)

// 添加执行方法，可以给定名字，以便于后续执行删除
func Add(spec string, f func(), name ... string) error {
    if len(name) > 0 {
        if Search(name[0]) != nil {
            return errors.New(fmt.Sprintf(`cron job "%s" already exists`, name[0]))
        }
        c := cron.New()
        if err := c.AddFunc(spec, f); err == nil {
            cronEntries.Append(Entry{
                Spec : spec,
                Cmd  : runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
                Time : gtime.Now(),
                Name : name[0],
                cron : c,
            })
            go c.Run()
        } else {
            return err
        }
    } else {
        if err := defaultCron.AddFunc(spec, f); err == nil {
            if cronStatus.Add(1) == 1 {
                go defaultCron.Run()
            }
            cronEntries.Append(Entry{
                Spec : spec,
                Cmd  : runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
                Time : gtime.Now(),
            })
        } else {
            return err
        }
    }
    return nil
}

// 检索指定名称的定时任务
func Search(name string) *Entry {
    entry, _ := searchEntry(name)
    return entry
}

// 检索指定名称的定时任务
func searchEntry(name string) (*Entry, int) {
    entry := (*Entry)(nil)
    index := -1
    cronEntries.RLockFunc(func(array []interface{}) {
        for k, v := range array {
            e := v.(Entry)
            if e.Name == name {
                entry = &e
                index = k
                break
            }
        }
    })
    return entry, index
}

// 根据指定名称删除定时任务
func Remove(name string) {
    if entry, index := searchEntry(name); index >= 0 {
        entry.cron.Stop()
        cronEntries.Remove(index)
    }
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

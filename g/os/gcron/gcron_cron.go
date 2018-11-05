// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron

import (
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/third/github.com/robfig/cron"
    "reflect"
    "runtime"
    "time"
)

// 添加定时任务
func (c *Cron) Add(spec string, f func(), name ... string) error {
    if len(name) > 0 {
        if Search(name[0]) != nil {
            return errors.New(fmt.Sprintf(`cron job "%s" already exists`, name[0]))
        }
        jobCron := cron.New()
        if err := jobCron.AddFunc(spec, f); err == nil {
            entry := &Entry{
                Spec   : spec,
                Cmd    : runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
                Time   : gtime.Now(),
                Name   : name[0],
                Status : gtype.NewInt(0),
                cron   : jobCron,
            }
            entry.Start()
            c.entries.Append(entry)
        } else {
            return err
        }
    } else {
        if err := c.cron.AddFunc(spec, f); err == nil {
            entry := &Entry {
                Spec   : spec,
                Cmd    : runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
                Time   : gtime.Now(),
                Status : c.status,
                cron   : c.cron,
            }
            entry.Start()
            c.entries.Append(entry)
        } else {
            return err
        }
    }
    return nil
}

// 延迟添加定时任务，delay参数单位为秒
func (c *Cron) DelayAdd(delay int, spec string, f func(), name ... string) {
    gtime.SetTimeout(time.Duration(delay)*time.Second, func() {
        if err := c.Add(spec, f, name ...); err != nil {
            panic(err)
        }
    })
}

// 检索指定名称的定时任务
func (c *Cron) Search(name string) *Entry {
    entry, _ := c.searchEntry(name)
    return entry
}

// 检索指定名称的定时任务
func (c *Cron) searchEntry(name string) (*Entry, int) {
    entry := (*Entry)(nil)
    index := -1
    c.entries.RLockFunc(func(array []interface{}) {
        for k, v := range array {
            e := v.(*Entry)
            if e.Name == name {
                entry = e
                index = k
                break
            }
        }
    })
    return entry, index
}

// 根据指定名称删除定时任务
func (c *Cron) Remove(name string) {
    if entry, index := c.searchEntry(name); index >= 0 {
        entry.cron.Stop()
        c.entries.Remove(index)
    }
}

// 开启定时任务执行(可以指定特定名称的一个或若干个定时任务)
func (c *Cron) Start(name...string) {
    if len(name) > 0 {
        for _, v := range name {
            if entry := c.Search(v); entry != nil {
                entry.Start()
            }
        }
    } else {
        c.entries.RLockFunc(func(array []interface{}) {
            for _, v := range array {
                v.(*Entry).Start()
            }
        })
    }
}

// 关闭定时任务执行(可以指定特定名称的一个或若干个定时任务)
func (c *Cron) Stop(name...string) {
    if len(name) > 0 {
        for _, v := range name {
            if entry := c.Search(v); entry != nil {
                entry.Stop()
            }
        }
    } else {
        c.entries.RLockFunc(func(array []interface{}) {
            for _, v := range array {
                v.(*Entry).Stop()
            }
        })
    }
}


// 获取所有已注册的定时任务项
func (c *Cron) Entries() []*Entry {
    length  := c.entries.Len()
    entries := make([]*Entry, length)
    for i := 0; i < length; i++ {
        entries[i] = c.entries.Get(i).(*Entry)
    }
    return entries
}

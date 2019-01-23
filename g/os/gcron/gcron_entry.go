// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron

import (
    "gitee.com/johng/gf/g/os/gtimer"
    "strconv"
    "time"
)

// 定时任务项
type Entry struct {
    cron       *Cron         // 所属定时任务
    entry      *gtimer.Entry // 定时器任务对象
    schedule   *cronSchedule // 定时任务配置对象
    Name       string        // 定时任务名称
    Job        func()        // 注册定时任务方法
    Time       time.Time     // 注册时间
}

// 创建定时任务
func (c *Cron) addEntry(pattern string, job func(), singleton bool, times int, name ... string) (*Entry, error) {
    schedule, err := newSchedule(pattern)
    if err != nil {
        return nil, err
    }
    entry := &Entry {
        cron      : c,
        schedule  : schedule,
        Job       : job,
        Time      : time.Now(),
    }
    if len(name) > 0 {
        entry.Name = name[0]
    } else {
        entry.Name = strconv.Itoa(c.idgen.Add(1))
    }
    entry.entry = gtimer.AddEntry(time.Second, entry.check, singleton, times, gtimer.STATUS_STOPPED)
    entry.entry.Start()
    c.entries.Set(entry.Name, entry)
    return entry, nil
}

// 是否单例运行
func (entry *Entry) IsSingleton() bool {
    return entry.entry.IsSingleton()
}

// 设置单例运行
func (entry *Entry) SetSingleton(enabled bool) {
    entry.entry.SetSingleton(true)
}

// 设置任务的运行次数
func (entry *Entry) SetTimes(times int) {
    entry.entry.SetTimes(times)
}

// 定时任务状态
func (entry *Entry) Status() int {
    return entry.entry.Status()
}

// 设置定时任务状态, 返回设置之前的状态
func (entry *Entry) SetStatus(status int) int {
    return entry.entry.SetStatus(status)
}

// 启动定时任务
func (entry *Entry) Start() {
    entry.entry.Start()
}

// 停止定时任务
func (entry *Entry) Stop() {
    entry.entry.Stop()
}

// 关闭定时任务
func (entry *Entry) Close() {
    entry.cron.Remove(entry.Name)
    entry.entry.Close()
}

// 定时任务检查执行
func (entry *Entry) check() {
    if entry.schedule.meet(time.Now()) {
        switch entry.cron.status.Val() {
            case STATUS_STOPPED:
                return

            case STATUS_CLOSED:
                entry.cron.Remove(entry.Name)
                gtimer.Exit()

            case STATUS_READY: fallthrough
            case STATUS_RUNNING:
                defer func() {
                    if entry.entry.Status() == STATUS_CLOSED {
                        entry.cron.Remove(entry.Name)
                    }
                }()
                entry.Job()
        }
    }
}

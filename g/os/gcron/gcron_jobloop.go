// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron

import (
    "gitee.com/johng/gf/g/container/garray"
    "time"
)

// 延迟添加定时任务，delay参数单位为秒
func (c *Cron) startLoop() {
    go func() {
        for c.status.Val() != STATUS_CLOSED {
            time.Sleep(time.Second)
            if c.status.Val() == STATUS_RUNNING {
                go c.checkEntries(time.Now())
            }
        }
    }()
}

// 遍历检查可执行定时任务，并异步执行
func (c *Cron) checkEntries(t time.Time) {
    removeArray := garray.NewStringArray(0, 0, false)
    c.entries.RLockFunc(func(m map[string]interface{}) {
        for _, v := range m {
            entry := v.(*Entry)
            if entry.schedule.meet(t) {
                // 是否已命令停止运行
                if entry.status.Val() == STATUS_CLOSED {
                    continue
                }
                switch entry.mode.Val() {
                    // 是否只允许单例运行
                    case MODE_SINGLETON:
                        if  entry.status.Set(STATUS_RUNNING) == STATUS_RUNNING {
                            continue
                        }
                    // 只运行一次的任务
                    case MODE_ONCE:
                        if  entry.status.Set(STATUS_CLOSED) == STATUS_CLOSED {
                            continue
                        }
                        removeArray.Append(entry.Name)
                }
                // 执行异步运行
                go func() {
                    defer func() {
                        if entry.status.Val() != STATUS_CLOSED {
                            entry.status.Set(STATUS_READY)
                        }
                    }()
                    entry.Job()
                }()
            }
        }
    })
    if removeArray.Len() > 0 {
        c.entries.BatchRemove(removeArray.Slice())
    }
}
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcron

// 启动定时任务
func (entry *Entry) Start() {
    if entry.Status.Set(1) == 0 {
        entry.cron.Start()
    }
}

// 关闭定时任务
func (entry *Entry) Stop() {
    if entry.Status.Set(0) == 1 {
        entry.cron.Stop()
    }
}

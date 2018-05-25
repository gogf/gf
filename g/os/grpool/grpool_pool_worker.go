// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool

import "gitee.com/johng/gf/g/os/gtime"

// 开始任务
func (w *PoolWorker) start() {
    go func() {
        for {
            if f := <- w.job; f != nil {
                // 执行任务
                f()
                // 更新活动时间(不存在并发安全问题)
                w.update = gtime.Second()
                // 执行完毕后添加到空闲队列
                if !w.pool.addWorker(w) {
                    break
                }
            } else {
                break
            }
        }
        w.pool.workerNum.Add(-1)
    }()
}

// 关闭当前任务
func (w *PoolWorker) stop() {
    w.setJob(nil)
}

// 设置当前任务的执行函数
func (w *PoolWorker) setJob(f func()) {
    w.job <- f
}


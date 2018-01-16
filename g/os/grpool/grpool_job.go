// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool

import "gitee.com/johng/gf/g/os/gtime"

// 开始任务
func (j *PoolJob) start() {
    go func() {
        for {
            if f := <- j.job; f != nil {
                // 执行任务
                f()
                // 更新活动时间
                j.update = gtime.Second()
                // 执行完毕后添加到空闲队列
                if !j.pool.addJob(j) {
                    break
                }
            } else {
                break
            }
        }
    }()
}

// 关闭当前任务
func (j *PoolJob) stop() {
    j.setJob(nil)
}

// 设置当前任务的执行函数
func (j *PoolJob) setJob(f func()) {
    j.job <- f
}


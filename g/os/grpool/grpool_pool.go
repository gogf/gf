// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool

import (
    "time"
    "sync/atomic"
    "gitee.com/johng/gf/g/os/gtime"
)

// 任务分配循环
func (p *Pool) startWorkLoop() {
    go func() {
        for {
            select {
                case <-p.funcEvents:
                    p.getJob().setJob(p.funcs.PopFront().(func()))
                case <-p.stopEvents:
                    return
            }
        }
    }()
}

// 定时清理过期任务
func (p *Pool) startClearLoop() {
    go func() {
        for {
            time.Sleep(gDEFAULT_CLEAR_INTERVAL*time.Second)
            if len(p.stopEvents) > 0 || len(p.funcEvents) == 0 {
                var j *PoolJob
                for {
                    if r := p.queue.PopFront(); r != nil {
                        j = r.(*PoolJob)
                        if gtime.Second() - int64(p.expire) > j.update {
                            j.stop()
                            atomic.AddInt32(&p.number, -1)
                        } else {
                            p.queue.PushFront(r)
                            break
                        }
                    } else {
                        break
                    }
                }
            }
            // 判断是池已经关闭，是则退出
            if len(p.stopEvents) > 0 {
                break
            }
        }
    }()
}

// 获取过期时间
func (p *Pool) getExpire() int32 {
    return atomic.LoadInt32(&p.expire)
}

// 创建一个空的任务对象
func (p *Pool) newJob() *PoolJob {
    j := &PoolJob {
        job  : make(chan func(), 1),
        pool : p,
    }
    j.start()
    atomic.AddInt32(&p.number, 1)
    return j
}

// 添加任务对象到队列
func (p *Pool) addJob(j *PoolJob) bool {
    if j.pool.getExpire() == -1 {
        return false
    }
    return p.queue.PushBack(j) != nil
}

// 获取/创建任务
func (p *Pool) getJob() *PoolJob {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolJob)
    }
    return p.newJob()
}

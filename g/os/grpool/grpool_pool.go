// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool

import (
    "time"
    "runtime"
    "sync/atomic"
    "gitee.com/johng/gf/g/os/gtime"
)

// 任务分配循环协程，使用基于 runtime.GOMAXPROCS 数量的协程来实现抢占调度
func (p *Pool) startWorkLoop() {
    for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
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
}

// 定时清理过期任务，单线程处理
func (p *Pool) startClearLoop() {
    go func() {
        for {
            // 如果接收到关闭通知(池已经关闭)，那么不再执行清理操作，直接退出
            if len(p.stopEvents) > 0 {
                break
            }
            time.Sleep(gDEFAULT_CLEAR_INTERVAL*time.Second)
            if len(p.funcEvents) == 0 {
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
        }
    }()
}

// 判断是否达到goroutine上上限
func (p *Pool) reachSizeLimit() bool {
    return atomic.LoadInt32(&p.number) >= atomic.LoadInt32(&p.size)
}

// 获取过期时间
func (p *Pool) getExpire() int32 {
    return atomic.LoadInt32(&p.expire)
}

// 创建一个空的任务对象
func (p *Pool) newJob() *PoolJob {
    // 如果达到goroutine数限制，那么阻塞等待有空闲goroutine后继续
    if p.reachSizeLimit() {
        // 阻塞等待空闲的协程资源，
        // 这是一个递归循环，因为该流程中存在协程抢占机制，
        // 如果进入getJob方法没有抢占到协程资源，那么该任务执行会继续等待下一个freeEvents
        <- p.freeEvents
        return p.getJob()
    }
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
    p.queue.PushBack(j)
    // 如果当前的goroutine数量达到上线，那么需要使用空闲goroutine通知事件
    if p.reachSizeLimit() {
        p.freeEvents <- struct{}{}
    }
    return true
}

// 获取/创建任务
func (p *Pool) getJob() *PoolJob {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolJob)
    }
    return p.newJob()
}

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
            time.Sleep(gDEFAULT_CLEAR_INTERVAL*time.Second)
            // 保证没有工作任务的情况下，执行worker清理操作
            if len(p.funcEvents) == 0 {
                var w *PoolWorker
                for {
                    if r := p.queue.PopFront(); r != nil {
                        w = r.(*PoolWorker)
                        if gtime.Second() - int64(p.expire) > w.update {
                            w.stop()
                            atomic.AddInt32(&p.number, -1)
                        } else {
                            p.queue.PushFront(w)
                            break
                        }
                    } else {
                        break
                    }
                }
            }
            // 如果接收到关闭通知(池已经关闭)，闭所有worker后退出
            if len(p.stopEvents) > 0 {
                for {
                    if r := p.queue.PopFront(); r != nil {
                        // 主动关闭所有worker，防止goroutine泄露
                        r.(*PoolWorker).stop()
                        atomic.AddInt32(&p.number, -1)
                    } else {
                        break
                    }
                }
                break
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
func (p *Pool) newJob() *PoolWorker {
    // 如果达到goroutine数限制，那么阻塞等待有空闲goroutine后继续
    if p.reachSizeLimit() {
        // 阻塞等待空闲的协程资源，
        // 这是一个递归循环，因为该流程中存在协程抢占机制，
        // 如果进入getJob方法没有抢占到协程资源，那么该任务执行会继续等待下一个freeEvents
        <- p.freeEvents
        return p.getJob()
    }
    w := &PoolWorker {
        job  : make(chan func(), 1),
        pool : p,
    }
    w.start()
    atomic.AddInt32(&p.number, 1)
    return w
}

// 添加任务对象到队列
func (p *Pool) addJob(w *PoolWorker) bool {
    if w.pool.getExpire() == -1 {
        return false
    }
    p.queue.PushBack(w)
    // 如果当前的goroutine数量达到上线，那么需要使用空闲goroutine通知事件
    if p.reachSizeLimit() {
        p.freeEvents <- struct{}{}
    }
    return true
}

// 获取/创建任务
func (p *Pool) getJob() *PoolWorker {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolWorker)
    }
    return p.newJob()
}

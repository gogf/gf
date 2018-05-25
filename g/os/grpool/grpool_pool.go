// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package grpool

import (
    "time"
    "runtime"
    "gitee.com/johng/gf/g/os/gtime"
)

// 任务分配循环协程，使用基于 runtime.GOMAXPROCS 数量的协程来实现抢占调度，
// 使用抢占调度的目的是使得任务能够并发地快速被分配出去执行
func (p *Pool) startSchedLoop() {
    for i := 0; i < runtime.GOMAXPROCS(-1); i++ {
        go func() {
            for {
                select {
                    case <-p.jobEvents:
                        p.getWorker().setJob(p.funcs.PopFront().(func()))
                    case <-p.stopEvents:
                        return
                }
            }
        }()
    }
}

// 定时清理过期任务，单协程处理
func (p *Pool) startClearLoop() {
    go func() {
        for {
            select {
                case <-p.stopEvents:
                    // 如果接收到关闭通知(池已经关闭)，关闭所有worker后退出
                    for {
                        if r := p.queue.PopFront(); r != nil {
                            // 主动关闭所有worker，防止goroutine泄露
                            r.(*PoolWorker).stop()
                        } else {
                            break
                        }
                    }
                    return

                default:
                    time.Sleep(gDEFAULT_CLEAR_INTERVAL*time.Second)
                    // 保证没有工作任务的情况下，执行worker清理操作
                    if len(p.jobEvents) == 0 {
                        var w *PoolWorker
                        for {
                            if r := p.queue.PopFront(); r != nil {
                                w = r.(*PoolWorker)
                                if gtime.Second() - int64(p.expire.Val()) > w.update {
                                    w.stop()
                                } else {
                                    p.queue.PushFront(w)
                                    break
                                }
                            } else {
                                break
                            }
                        }
                    }
            }
        }
    }()
}

// 获取过期时间
func (p *Pool) getExpire() int {
    return p.expire.Val()
}

// 创建一个空的任务对象
func (p *Pool) newWorker() *PoolWorker {
    // 如果达到goroutine数限制，那么阻塞等待有空闲goroutine后继续
    // 需要注意的是在高并发下workerNum的值可能会高于size，
    // 从效率上考虑没有将workerNum和size都放到一个互斥锁中进行准确度控制，
    // 精准是要付出代价的
    if p.workerNum.Val() >= p.size.Val() {
        // (非精准控制)阻塞等待空闲的协程资源，
        // 这是一个递归循环，因为该流程中存在协程抢占机制，
        // 如果进入getJob方法没有抢占到协程资源，那么该任务执行会继续等待下一个freeEvents事件产生
        p.blockedNum.Add(1)
        <- p.freeEvents
        return p.getWorker()
    }
    w := &PoolWorker {
        job  : make(chan func(), 1),
        pool : p,
    }
    w.start()
    p.workerNum.Add(1)
    return w
}

// 添加worker对象到空闲队列
func (p *Pool) addWorker(w *PoolWorker) bool {
    if p.workerNum.Val() > p.size.Val() || w.pool.getExpire() == -1 {
        return false
    }
    p.queue.PushBack(w)
    // 如果当前的goroutine数量达到上线，那么需要使用空闲goroutine通知事件
    if p.blockedNum.Val() > 0 {
        p.blockedNum.Add(-1)
        p.freeEvents <- struct{}{}
    }
    return true
}

// 获取/创建任务
func (p *Pool) getWorker() *PoolWorker {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolWorker)
    }
    return p.newWorker()
}

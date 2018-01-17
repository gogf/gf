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

// 任务分配循环，这是一个独立的goroutine，单线程处理
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

// 定时清理过期任务，单线程处理
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
            // 判断是池已经关闭，并且所有goroutine已退出，那么该goroutine终止执行
            if len(p.stopEvents) > 0 && atomic.LoadInt32(&p.number) == 0 {
                break
            }
        }
    }()
}

// 判断是否达到goroutine上上限
func (p *Pool) reachSizeLimit() bool {
    return false
    return atomic.LoadInt32(&p.number) >= atomic.LoadInt32(&p.size)
}

// 获取过期时间
func (p *Pool) getExpire() int32 {
    return atomic.LoadInt32(&p.expire)
}

// 创建一个空的任务对象
func (p *Pool) newJob() *PoolJob {
    // 如果达到goroutine数限制，那么阻塞等待有空闲goroutine后继续
    //if p.reachSizeLimit() {
    //    // 阻塞等待空闲goroutine
    //    select {
    //        case <- p.freeEvents:
    //            return p.getJob()
    //    }
    //}
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
    //if p.reachSizeLimit() {
    //    p.freeEvents <- struct{}{}
    //}
    return true
}

// 获取/创建任务
func (p *Pool) getJob() *PoolJob {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolJob)
    }
    return p.newJob()
}

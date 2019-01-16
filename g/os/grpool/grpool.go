// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package grpool implements a goroutine reusable pool.
// 
// Goroutine池,
// 用于goroutine复用，提升异步操作执行效率(避免goroutine限制，并节约内存开销).
// 需要注意的是，grpool提供给的公共池不提供关闭方法，自创建的池可以手动关闭掉。
package grpool

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
    "math"
)

// goroutine池对象
type Pool struct {
    workerChan  chan struct{}      // 使用channel限制最大的goroutine数量
    workerNum   *gtype.Int         // 当前正在运行的worker/goroutine数量
    jobQueue    *glist.List        // 待处理任务操作队列
    jobEvents   chan struct{}      // 任务添加事件(jobQueue+jobEvents结合使用)
    closed      *gtype.Bool
}

// 默认的goroutine池管理对象
// 该对象与进程同生命周期，无需Close
var defaultPool = New()

// 创建goroutine池管理对象， 参数用于限制限制最大的goroutine数量/线程数/worker数量，非必需参数，默认不做限制
func New(size...int) *Pool {
    s := 0
    if len(size) > 0 {
        s = size[0]
    }
    p := &Pool {
        workerNum   : gtype.NewInt(),
        jobQueue    : glist.New(),
        jobEvents   : make(chan struct{}, math.MaxInt32),
        workerChan  : make(chan struct{}, s),
        closed      : gtype.NewBool(),
    }
    return p
}

// 添加异步任务(使用默认的池对象)
func Add(f func()) error {
    return defaultPool.Add(f)
}

// 查询当前goroutine总数
func Size() int {
    return defaultPool.workerNum.Val()
}

// 查询当前等待处理的任务总数
func Jobs() int {
    return len(defaultPool.jobEvents)
}

// 添加异步任务
func (p *Pool) Add(f func()) error {
    p.jobQueue.PushBack(f)
    p.jobEvents <- struct{}{}
    // 判断是否创建新的worker
    if p.Jobs() > 1 || p.workerNum.Val() == 0 {
        p.ForkWorker()
    }
    return nil
}

// 查询当前goroutine worker总数
func (p *Pool) Size() int {
    return p.workerNum.Val()
}

// 查询当前等待处理的任务总数
func (p *Pool) Jobs() int {
    return p.jobQueue.Len()
}

// 创建新的worker执行任务
func (p *Pool) ForkWorker() {
    if cap(p.workerChan) > 0 {
        // 如果worker数量已经达到限制，那么不创建新worker，直接返回
        if p.workerNum.Val() == cap(p.workerChan) {
            return
        }
        p.workerNum.Add(1)
        p.workerChan <- struct{}{}
    } else {
        p.workerNum.Add(1)
    }
    go func() {
        for !p.closed.Val() {
            select {
                case <- p.jobEvents:
                    if job := p.jobQueue.PopFront(); job != nil {
                        job.(func())()
                    } else {
                        goto WorkerDone
                    }
                default:
                    goto WorkerDone
            }
        }
WorkerDone:
        p.workerNum.Add(-1)
        if cap(p.workerChan) > 0 {
            <- p.workerChan
        }
    }()
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    p.closed.Set(true)
}
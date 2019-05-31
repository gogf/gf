// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
// 
// Goroutine池,
// 用于goroutine复用，提升异步操作执行效率(避免goroutine限制，并节约内存开销).
// 需要注意的是，grpool提供给的公共池不提供关闭方法，自创建的池可以手动关闭掉。
package grpool

import (
    "github.com/gogf/gf/g/container/glist"
    "github.com/gogf/gf/g/container/gtype"
    "math"
)

// goroutine池对象
type Pool struct {
    limit       int                // 最大的goroutine数量限制
    count       *gtype.Int         // 当前正在运行的goroutine数量
    list        *glist.List        // 待处理任务操作列表
    events      chan struct{}      // 任务添加事件
    closed      *gtype.Bool
}

// 默认的goroutine池管理对象
// 该对象与进程同生命周期，无需Close
var defaultPool = New()

// 创建goroutine池管理对象，参数用于限制限制最大的goroutine数量，非必需参数，默认不做限制
func New(size...int) *Pool {
    limit := -1
    if len(size) > 0 {
	    limit = size[0]
    }
    p := &Pool {
	    limit  : limit,
        count  : gtype.NewInt(),
        list   : glist.New(),
        events : make(chan struct{}, math.MaxInt32),
        closed : gtype.NewBool(),
    }
    return p
}

// 添加异步任务(使用默认的池对象)
func Add(f func()) error {
    return defaultPool.Add(f)
}

// 查询当前goroutine总数
func Size() int {
    return defaultPool.count.Val()
}

// 查询当前等待处理的任务总数
func Jobs() int {
    return len(defaultPool.events)
}

// 添加异步任务
func (p *Pool) Add(f func()) error {
    p.list.PushBack(f)
    p.events <- struct{}{}
    // 判断是否创建新的worker
    if p.list.Len() > 1 || p.count.Val() == 0 {
        p.fork()
    }
    return nil
}

// 查询当前goroutine总数
func (p *Pool) Size() int {
    return p.count.Val()
}

// 查询当前等待处理的任务总数
func (p *Pool) Jobs() int {
    return p.list.Len()
}

// 创建新的worker执行任务
func (p *Pool) fork() {
    // 如果worker数量已经达到限制，那么不创建新worker，直接返回
    if p.count.Val() == p.limit {
        return
    }
	p.count.Add(1)
    go func() {
        for !p.closed.Val() {
            select {
                case <- p.events:
                    if job := p.list.PopFront(); job != nil {
                        job.(func())()
                    } else {
	                    p.count.Add(-1)
	                    return
                    }
                default:
	                p.count.Add(-1)
	                return
            }
        }
    }()
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    p.closed.Set(true)
}
// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpool implements a goroutine reusable pool.
package grpool

import (
	"github.com/gogf/gf/g/container/glist"
	"github.com/gogf/gf/g/container/gtype"
)

// Goroutine Pool
type Pool struct {
    limit  int         // 最大的goroutine数量限制
    count  *gtype.Int  // 当前正在运行的goroutine数量
    list   *glist.List // 待处理任务操作列表
    closed *gtype.Bool // 是否关闭
}

// 默认的goroutine池管理对象,
// 该对象与进程同生命周期，无需Close
var defaultPool = New()

// 创建goroutine池管理对象，参数用于限制限制最大的goroutine数量，非必需参数，默认不做限制
func New(limit...int) *Pool {
    p := &Pool {
	    limit  : -1,
        count  : gtype.NewInt(),
        list   : glist.New(),
        closed : gtype.NewBool(),
    }
    if len(limit) > 0 {
    	p.limit = limit[0]
    }
    return p
}

// 添加异步任务(使用默认的池对象)
func Add(f func()) {
    defaultPool.Add(f)
}

// 查询当前goroutine总数
func Size() int {
    return defaultPool.count.Val()
}

// 查询当前等待处理的任务总数
func Jobs() int {
    return defaultPool.list.Len()
}

// 添加异步任务
func (p *Pool) Add(f func()) {
    p.list.PushFront(f)
    // 判断是否创建新的goroutine
    if p.count.Val() != p.limit {
        p.fork()
    }
}

// 查询当前goroutine总数
func (p *Pool) Size() int {
    return p.count.Val()
}

// 查询当前等待处理的任务总数
func (p *Pool) Jobs() int {
    return p.list.Size()
}

// 检查并创建新的goroutine执行任务
func (p *Pool) fork() {
	p.count.Add(1)
    go func() {
    	defer p.count.Add(-1)
    	job := (interface{})(nil)
        for !p.closed.Val() {
        	if job = p.list.PopBack(); job != nil {
		        job.(func())()
	        } else {
	        	return
	        }
        }
    }()
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
	p.closed.Set(true)
}
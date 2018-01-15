// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Goroutine池.
// 用于goroutine复用，提升异步操作执行效率.
package groutine

import (
    "math"
    "sync"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/container/glist"
)

// goroutine池对象
type Pool struct {
    jobs   *gset.InterfaceSet // 当前任务对象(*PoolJob)
    queue  *glist.SafeList    // 空闲任务队列(*PoolJob)
    funcs  *glist.SafeList    // 待处理任务操作队列
    events chan struct{}      // 任务操作处理事件(用于任务事件通知)
}

// goroutine任务
type PoolJob struct {
    mu     sync.RWMutex
    job    chan func() // 当前任务(当为nil时表示关闭)
    pool   *Pool       // 所属池
}

// 创建goroutine池管理对象
func New() *Pool {
    p := &Pool {
        jobs   : gset.NewInterfaceSet(),
        queue  : glist.NewSafeList(),
        funcs  : glist.NewSafeList(),
        events : make(chan struct{}, math.MaxUint32),
    }
    p.loop()
    return p
}

// 添加异步任务
func (p *Pool) Add(f func()) {
    p.funcs.PushBack(f)
    p.events <- struct{}{}
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    p.Add(nil)
    p.jobs.Iterator(func(v interface{}){
        v.(*PoolJob).stop()
    })
}
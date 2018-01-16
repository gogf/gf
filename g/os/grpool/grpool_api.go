// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Goroutine池.
// 用于goroutine复用，提升异步操作执行效率.
package grpool

import (
    "math"
    "sync/atomic"
    "gitee.com/johng/gf/g/container/glist"
)

const (
    gDEFAULT_EXPIRE_TIME    = 60 // 默认goroutine过期时间
    gDEFAULT_CLEAR_INTERVAL = 60 // 定期检查任务过期时间间隔
)

// goroutine池对象
type Pool struct {
    expire     int32              // goroutine过期时间(秒)
    number     int32              // 当前goroutine数量(非任务数)
    queue      *glist.SafeList    // 空闲任务队列(*PoolJob)
    funcs      *glist.SafeList    // 待处理任务操作队列
    funcEvents chan struct{}      // 任务操作处理事件(用于任务事件通知)
    stopEvents chan struct{}      // 池关闭事件(用于池相关异步线程通知)
}

// goroutine任务
type PoolJob struct {
    job    chan func() // 当前任务(当为nil时表示关闭)
    pool   *Pool       // 所属池
    update int64       // 更新时间
}

// 默认的goroutine池管理对象
// 该对象与进程同生命周期，无需Close
var defaultPool = New(gDEFAULT_EXPIRE_TIME)

// 创建goroutine池管理对象，给定过期时间(秒)
func New(expire int) *Pool {
    p := &Pool {
        expire     : int32(expire),
        queue      : glist.NewSafeList(),
        funcs      : glist.NewSafeList(),
        funcEvents : make(chan struct{}, math.MaxUint32),
        stopEvents : make(chan struct{}, 1),
    }
    p.startWorkLoop()
    p.startClearLoop()
    return p
}

// 添加异步任务(使用默认的池对象)
func Add(f func()) {
    defaultPool.funcs.PushBack(f)
    defaultPool.funcEvents <- struct{}{}
}

// 查询当前goroutine总数
func Size() int {
    return int(atomic.LoadInt32(&defaultPool.number))
}

// 设置默认池中goroutine的过期时间
func SetExpire(expire int) {
    atomic.StoreInt32(&defaultPool.expire, int32(expire))
}

// 添加异步任务
func (p *Pool) Add(f func()) {
    p.funcs.PushBack(f)
    p.funcEvents <- struct{}{}
}

// 查询当前goroutine总数
func (p *Pool) Size() int {
    return int(atomic.LoadInt32(&p.number))
}

// 设置当前池中goroutine的过期时间
func (p *Pool) SetExpire(expire int) {
    atomic.StoreInt32(&p.expire, int32(expire))
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    // 必须首先标识让任务过期自动关闭
    p.SetExpire(-1)
    p.stopEvents <- struct{}{} // 通知workloop
    p.stopEvents <- struct{}{} // 通知clearloop

}
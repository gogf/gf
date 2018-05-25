// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Goroutine池.
// 用于goroutine复用，提升异步操作执行效率.
// 需要注意的是，grpool提供给的公共池不提供关闭方法(但可以修改公共属性)，自创建的池可以手动关闭掉。
package grpool

import (
    "math"
    "errors"
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gtype"
)

const (
    gDEFAULT_EXPIRE_TIME    = 60 // 默认goroutine过期时间(秒)
    gDEFAULT_CLEAR_INTERVAL = 60 // 定期检查任务过期时间间隔(秒)
)

// goroutine池对象
type Pool struct {
    size        *gtype.Int         // 限制最大的goroutine数量/协程数/worker数量
    expire      *gtype.Int         // goroutine过期时间(秒)
    workerNum   *gtype.Int         // 当前正在运行的goroutine数量(非任务数)
    blockedNum  *gtype.Int         // 当前被阻塞运行的goroutine数量
    queue       *glist.List        // 空闲任务队列(*PoolJob)
    funcs       *glist.List        // 待处理任务操作队列
    jobEvents   chan struct{}      // 任务添加事件(兄弟们该干活了！)
    freeEvents  chan struct{}      // 空闲协程通知事件
    stopEvents  chan struct{}      // 池关闭事件(用于池相关异步协程通知)
}

// 一个worker对应一个goroutine
type PoolWorker struct {
    job    chan func() // 当前任务(当为nil时表示关闭)
    pool   *Pool       // 所属协程池
    update int64       // 更新时间
}

// 默认的goroutine池管理对象
// 该对象与进程同生命周期，无需Close
var defaultPool = New(gDEFAULT_EXPIRE_TIME)

// 创建goroutine池管理对象，给定过期时间(秒)
// 第二个参数用于限制限制最大的goroutine数量/线程数/worker数量，非必需参数，默认不做限制
func New(expire int, size...int) *Pool {
    s := math.MaxInt32
    if len(size) > 0 {
        s = size[0]
    }
    p := &Pool {
        size        : gtype.NewInt(s),
        expire      : gtype.NewInt(expire),
        workerNum   : gtype.NewInt(),
        blockedNum  : gtype.NewInt(),
        queue       : glist.New(),
        funcs       : glist.New(),
        jobEvents   : make(chan struct{}, math.MaxInt32),
        freeEvents  : make(chan struct{}, math.MaxInt32),
        stopEvents  : make(chan struct{}, 0),
    }
    p.startSchedLoop()
    p.startClearLoop()
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

// 动态改变默认池中goroutine的上线数量
func SetSize(size int) {
    defaultPool.SetSize(size)
}

// 动态改变默认池中goroutine的过期时间
func SetExpire(expire int) {
    defaultPool.SetExpire(expire)
}

// 添加异步任务
func (p *Pool) Add(f func()) error {
    if len(p.stopEvents) > 0 {
        return errors.New("pool closed")
    }
    p.funcs.PushBack(f)
    p.jobEvents <- struct{}{}
    return nil
}

// 查询当前goroutine worker总数
func (p *Pool) Size() int {
    return p.workerNum.Val()
}

// 查询当前等待处理的任务总数
func (p *Pool) Jobs() int {
    return len(p.jobEvents)
}

// 动态改变当前池中goroutine的上线数量
func (p *Pool) SetSize(size int) {
    p.size.Set(size)
}

// 动态改变当前池中goroutine的过期时间
func (p *Pool) SetExpire(expire int) {
    p.expire.Set(expire)
}

// 关闭池，所有的任务将会停止，此后继续添加的任务将不会被执行
func (p *Pool) Close() {
    // 必须首先标识让任务过期自动关闭
    p.SetExpire(-1)
    // 使用stopEvents事件通知所有的异步协程及清理协程自动退出
    close(p.stopEvents)
}
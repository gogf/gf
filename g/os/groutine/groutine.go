// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Goroutine池.
package groutine

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gset"
    "sync"
)

// goroutine池对象
type Pool struct {
    queue *glist.SafeList    // 空闲任务队列*PoolJob)
    pjobs *gset.InterfaceSet // 当前任务对象(*PoolJob)
}

// goroutine任务
type PoolJob struct {
    mu     sync.RWMutex
    job    chan func() // 当前任务(当为nil时表示关闭)
    pool   *Pool       // 所属池
}

// 创建一个空的任务对象
func (p *Pool) newJob() *PoolJob {
    j := &PoolJob {
        job  : make(chan func(), 1),
        pool : p,
    }
    j.start()
    p.pjobs.Add(j)
    return j
}

// 添加任务对象到队列
func (p *Pool) addJob(j *PoolJob) {
    p.queue.PushBack(j)
}

// 获取/创建任务
func (p *Pool) getJob() *PoolJob {
    if r := p.queue.PopFront(); r != nil {
        return r.(*PoolJob)
    }
    return p.newJob()
}

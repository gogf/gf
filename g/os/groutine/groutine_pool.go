// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package groutine

// 任务分配循环
func (p *Pool) loop() {
    go func() {
        for {
            // 阻塞监听任务事件
            if _, ok := <- p.events; ok {
                // 如果任务为nil，表示池关闭
                if r := p.funcs.PopFront(); r != nil {
                    p.getJob().setJob(r.(func()))
                } else {
                    return
                }
            }
        }
    }()
}

// 创建一个空的任务对象
func (p *Pool) newJob() *PoolJob {
    j := &PoolJob {
        job  : make(chan func(), 1),
        pool : p,
    }
    j.start()
    p.jobs.Add(j)
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

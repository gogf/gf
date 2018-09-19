// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 并发安全的动态队列.
// 特点：
// 1、动态队列初始化速度快；
// 2、动态的队列大小(不限大小)；
// 3、取数据时如果队列为空那么会阻塞等待；
package gqueue

import (
    "gitee.com/johng/gf/g/container/glist"
    "math"
)

// 0、这是一个先进先出的队列(chan <-- list)；
// 1、当创建Queue对象时限定大小，那么等同于一个同步的chan并发安全队列；
// 2、不限制大小时，list链表用以存储数据，临时chan负责为客户端读取数据，当从chan获取数据时，list往chan中不停补充数据；
// 3、由于功能主体是chan，那么操作仍然像chan那样具有阻塞效果；
type Queue struct {
    limit     int              // 队列限制大小
    queue     chan interface{} // 用于队列写入限制
    list      *glist.List      // 数据链表
    events    chan struct{}    // 通知chan，当不限制队列大小时的写入事件通知
    closeChan chan struct{}    // 关闭channel
}

const (
    // 默认临时队列大小,注意是临时的
    gDEFAULT_QUEUE_SIZE = 10000
)

// 队列大小为非必须参数，默认不限制
func New(limit...int) *Queue {
    size := gDEFAULT_QUEUE_SIZE
    if len(limit) > 0 {
        size = limit[0]
    }
    q := &Queue {
        list      : glist.New(),
        queue     : make(chan interface{}, size),
        events    : make(chan struct{}, math.MaxInt32),
        closeChan : make(chan struct{}, 0),
    }
    if len(limit) > 0 {
        q.limit = size
    } else {
        // 如果是动态队列大小，那么额外会运行一个goroutine
        go q.startAsyncLoop()
    }
    return q
}

// 异步list->chan同步队列
func (q *Queue) startAsyncLoop() {
    for {
        select {
            case <- q.closeChan:
                return
            case <- q.events:
                // 循环读取链表，直到为空才跳出
                for {
                    if v := q.list.PopFront(); v != nil {
                        q.queue <- v
                    } else {
                        break
                    }
                }
        }
    }
}

// 将数据压入队列, 队头
func (q *Queue) Push(v interface{}) {
    if q.limit > 0 {
        q.queue <- v
    } else {
        q.list.PushBack(v)
        if len(q.events) == 0 {
            q.events <- struct{}{}
        }
    }
}

// 从队头先进先出地从队列取出一项数据
func (q *Queue) Pop() interface{} {
    return <- q.queue
}

// 关闭队列(通知所有通过Pop*阻塞的协程退出)
func (q *Queue) Close() {
    q.list.RemoveAll()
    close(q.queue)
    close(q.events)
    close(q.closeChan)
}

// 获取当前队列大小
func (q *Queue) Size() int {
    return len(q.queue) + q.list.Len()
}



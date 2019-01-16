// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gqueue provides a dynamic/static concurrent-safe(alternative) queue.
//
// 并发安全的动态队列.
//
//   特点：
//   1. 动态队列初始化速度快；
//   2. 动态的队列大小(不限大小)；
//   3. 取数据时如果队列为空那么会阻塞等待；
package gqueue

import (
    "container/list"
    "math"
    "sync"
)

// 1、这是一个先进先出的队列(chan <-- list)；
//
// 2、当创建Queue对象时限定大小，那么等同于一个同步的chan并发安全队列；
//
// 3、不限制大小时，list链表用以存储数据，临时chan负责为客户端读取数据，当从chan获取数据时，list往chan中不停补充数据；
//
// 4、由于功能主体是chan，那么操作仍然像chan那样具有阻塞效果；
type Queue struct {
    mu        sync.Mutex       // 底层链表写锁
    limit     int              // 队列限制大小
    list      *list.List       // 底层数据链表
    events    chan struct{}    // 写入事件通知
    closed    chan struct{}    // 队列关闭通知
    C         chan interface{} // 队列数据读取
}

const (
    // 动态队列缓冲区大小
    gDEFAULT_QUEUE_SIZE = 10000
)

// 队列大小为非必须参数，默认不限制
func New(limit...int) *Queue {
    q := &Queue {
        closed : make(chan struct{}, 0),
    }
    if len(limit) > 0 {
        q.limit  = limit[0]
        q.C      = make(chan interface{}, limit[0])
    } else {
        q.list   = list.New()
        q.events = make(chan struct{}, math.MaxInt32)
        q.C      = make(chan interface{}, gDEFAULT_QUEUE_SIZE)
        go q.startAsyncLoop()
    }
    return q
}

// 异步list->chan同步队列
func (q *Queue) startAsyncLoop() {
    for {
        select {
            case <- q.closed:
                return
            case <- q.events:
                for {
                    if length := q.list.Len(); length > 0 {
                        array := make([]interface{}, length)
                        q.mu.Lock()
                        for i := 0; i < length; i++ {
                            if e := q.list.Front(); e != nil {
                                array[i] = q.list.Remove(e)
                            } else {
                                break
                            }
                        }
                        q.mu.Unlock()
                        for _, v := range array {
                           q.C <- v
                        }
                    } else {
                        break
                    }
                }
        }
    }
}

// 将数据压入队列, 队尾
func (q *Queue) Push(v interface{}) {
    if q.limit > 0 {
        q.C <- v
    } else {
        q.mu.Lock()
        q.list.PushBack(v)
        q.mu.Unlock()
        q.events <- struct{}{}
    }
}

// 从队头先进先出地从队列取出一项数据
func (q *Queue) Pop() interface{} {
    return <- q.C
}

// 关闭队列(通知所有通过Pop*阻塞的协程退出)
func (q *Queue) Close() {
    close(q.C)
    close(q.events)
    close(q.closed)
}

// 获取当前队列大小
func (q *Queue) Size() int {
    return len(q.C) + q.list.Len()
}



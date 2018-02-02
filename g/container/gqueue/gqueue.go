// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 动态大小的安全队列(dynamic channel).
package gqueue

import (
    "math"
    "sync"
    "container/list"
    "sync/atomic"
    "errors"
)

type Queue struct {
    mu     sync.RWMutex
    list   *list.List     // 数据队列
    limit  int            // 队列限制大小
    limits chan struct{}  // 用于队列大小限制
    events chan struct{}  // 用于内部数据写入事件通知
    closed int32          // 队列是否关闭
}

// 队列大小为非必须参数，默认不限制
func New(limit...int) *Queue {
    size := 0
    if len(limit) > 0 {
        size = limit[0]
    }
    return &Queue {
        list   : list.New(),
        limits : make(chan struct{}, size),
        events : make(chan struct{}, math.MaxInt64),
    }
}

// 将数据压入队列, 队尾
func (q *Queue) PushBack(v interface{}) error {
    if atomic.LoadInt32(&q.closed) > 0 {
        return errors.New("queue closed")
    }
    if q.limit > 0 {
        q.limits <- struct{}{}
    }
    q.mu.Lock()
    q.list.PushBack(v)
    q.mu.Unlock()
    q.events <- struct{}{}
    return nil
}

// 将数据压入队列, 队头
func (q *Queue) PushFront(v interface{}) error {
    if atomic.LoadInt32(&q.closed) > 0 {
        return errors.New("queue closed")
    }
    if q.limit > 0 {
        q.limits <- struct{}{}
    }
    q.mu.Lock()
    q.list.PushFront(v)
    q.mu.Unlock()
    q.events <- struct{}{}
    return nil
}

// 从队头先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *Queue) PopFront() interface{} {
    select {
        case <- q.events:
            if q.limit > 0 {
                <- q.limits
            }
            q.mu.Lock()
            if elem := q.list.Front(); elem != nil {
                item := q.list.Remove(elem)
                q.mu.Unlock()
                return item
            }
            q.mu.Unlock()
    }
    return nil
}

// 从队尾先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
// 第二个返回值表示队列是否关闭
func (q *Queue) PopBack() interface{} {
    select {
        case <- q.events:
            if q.limit > 0 {
                <- q.limits
            }
            q.mu.Lock()
            if elem := q.list.Front(); elem != nil {
                item := q.list.Remove(elem)
                q.mu.Unlock()
                return item
            }
            q.mu.Unlock()
    }
    return nil
}

// 关闭队列(通知所有通过Pop阻塞的协程退出)
func (q *Queue) Close() {
    if atomic.LoadInt32(&q.closed) == 0 {
        atomic.StoreInt32(&q.closed, 1)
        close(q.limits)
        close(q.events)
    }
}

// 获取当前队列大小
func (q *Queue) Size() int {
    return len(q.events)
}
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 并发安全的动态队列.
// 优点：
// 1、队列初始化速度快；
// 2、可以向队头/队尾进行Push/Pop操作；
// 3、取数据时如果队列为空那么会阻塞等待；
package gqueue

import (
    "math"
    "sync"
    "errors"
    "container/list"
    "gitee.com/johng/gf/g/container/gtype"
)

type Queue struct {
    mu     sync.RWMutex   // 用于队列并发安全处理
    list   *list.List     // 数据队列
    limit  int            // 队列限制大小
    limits chan struct{}  // 用于队列写入限制
    events chan struct{}  // 用于队列出列限制
    closed *gtype.Bool    // 队列是否关闭
}

// 队列大小为非必须参数，默认不限制
func New(limit...int) *Queue {
    size := 0
    if len(limit) > 0 {
        size = limit[0]
    }
    return &Queue {
        list   : list.New(),
        limit  : size,
        limits : make(chan struct{}, size),
        events : make(chan struct{}, math.MaxInt32),
        closed : gtype.NewBool(),
    }
}

// 将数据压入队列, 队尾
func (q *Queue) PushBack(v interface{}) error {
    if q.closed.Val() {
        return errors.New("closed")
    }
    if q.limit > 0 {
        q.limits <- struct{}{}
    }
    q.mu.Lock()
    q.list.PushBack(v)
    q.mu.Unlock()
    if q.limit == 0 {
        q.events <- struct{}{}
    }
    return nil
}

// 将数据压入队列, 队头
func (q *Queue) PushFront(v interface{}) error {
    if q.closed.Val() {
        return errors.New("closed")
    }
    // 限制队列大小，使用channel进行阻塞限制
    if q.limit > 0 {
        q.limits <- struct{}{}
    }
    q.mu.Lock()
    q.list.PushFront(v)
    q.mu.Unlock()
    if q.limit == 0 {
        q.events <- struct{}{}
    }
    return nil
}

// 从队头先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *Queue) PopFront() interface{} {
    if q.closed.Val() {
        return nil
    }
    if q.limit > 0 {
        <- q.limits
    } else {
        <- q.events
    }
    q.mu.Lock()
    if elem := q.list.Front(); elem != nil {
        item := q.list.Remove(elem)
        q.mu.Unlock()
        return item
    }
    q.mu.Unlock()
    return nil
}

// 从队尾先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *Queue) PopBack() interface{} {
    if q.closed.Val() {
        return nil
    }
    if q.limit > 0 {
        <- q.limits
    } else {
        <- q.events
    }
    q.mu.Lock()
    if elem := q.list.Front(); elem != nil {
        item := q.list.Remove(elem)
        q.mu.Unlock()
        return item
    }
    q.mu.Unlock()
    return nil
}

// 关闭队列(通知所有通过Pop*阻塞的协程退出)
func (q *Queue) Close() {
    if !q.closed.Val() {
        q.closed.Set(true)
        close(q.limits)
        close(q.events)
    }
}

// 获取当前队列大小
func (q *Queue) Size() int {
    return len(q.events)
}



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
)

type Queue struct {
    mu     sync.RWMutex
    list   *list.List
    events chan struct{}
}

func New() *Queue {
    return &Queue {
        list   : list.New(),
        events : make(chan struct{}, math.MaxInt64),
    }
}

// 将数据压入队列, 队尾
func (q *Queue) PushBack(v interface{}) {
    q.mu.Lock()
    q.list.PushBack(v)
    q.mu.Unlock()
    q.events <- struct{}{}
}

// 将数据压入队列, 队头
func (q *Queue) PushFront(v interface{}) {
    q.mu.Lock()
    q.list.PushFront(v)
    q.mu.Unlock()
    q.events <- struct{}{}
}

// 从队头先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
// 第二个返回值表示队列是否关闭
func (q *Queue) PopFront() (interface{}, bool) {
    select {
        case <- q.events:
            q.mu.Lock()
            if elem := q.list.Front(); elem != nil {
                item := q.list.Remove(elem)
                q.mu.Unlock()
                return item, true
            }
            q.mu.Unlock()
    }
    return nil, false
}

// 从队尾先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
// 第二个返回值表示队列是否关闭
func (q *Queue) PopBack() (interface{}, bool)  {
    select {
    case <- q.events:
        q.mu.Lock()
        if elem := q.list.Front(); elem != nil {
            item := q.list.Remove(elem)
            q.mu.Unlock()
            return item, true
        }
        q.mu.Unlock()
    }
    return nil, false
}

// 关闭队列(通知所有通过Pop阻塞的协程退出)
func (q *Queue) Close() {
    q.events <- struct{}{}
}

// 获取当前队列大小
func (q *Queue) Size() int {
    return len(q.events)
}
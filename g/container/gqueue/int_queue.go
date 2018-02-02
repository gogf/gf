// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gqueue

import (
    "math"
    "sync"
    "container/list"
)

type IntQueue struct {
    mu     sync.RWMutex
    list   *list.List
    events chan struct{}
}

func NewIntQueue() *IntQueue {
    return &IntQueue{
        list   : list.New(),
        events : make(chan struct{}, math.MaxInt64),
    }
}

// 将数据压入队列
func (q *IntQueue) Push(v int) {
    q.mu.Lock()
    q.list.PushBack(v)
    q.mu.Unlock()
    q.events <- struct{}{}
}

// 先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *IntQueue) Pop() int {
    select {
        case <- q.events:
            q.mu.Lock()
            if elem := q.list.Front(); elem != nil {
                item := q.list.Remove(elem).(int)
                q.mu.Unlock()
                return item
            }
            q.mu.Unlock()
    }
    return 0
}

// 关闭队列(通知所有通过Pop阻塞的协程退出)
func (q *IntQueue) Close() {
    q.events <- struct{}{}
}

// 获取当前队列大小
func (q *IntQueue) Size() int {
    return len(q.events)
}
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

type InterfaceQueue struct {
    mu     sync.RWMutex
    list   *list.List
    events chan struct{}
}

func NewInterfaceQueue() *InterfaceQueue {
    return &InterfaceQueue {
        list   : list.New(),
        events : make(chan struct{}, math.MaxInt64),
    }
}

// 将数据压入队列
func (q *InterfaceQueue) Push(v interface{}) {
    q.mu.Lock()
    q.list.PushBack(v)
    q.mu.Unlock()
    q.events <- struct{}{}
}

// 先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *InterfaceQueue) Pop() interface{} {
    select {
        case <- q.events:
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
func (q *InterfaceQueue) Close() {
    q.events <- struct{}{}
}

// 获取当前队列大小
func (q *InterfaceQueue) Size() int {
    return len(q.events)
}
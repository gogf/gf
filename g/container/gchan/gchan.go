// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 优雅的channel操作.
package gchan

import (
    "sync"
    "errors"
    "sync/atomic"
)

type Chan struct {
    mu     sync.RWMutex
    list   chan interface{}
    closed int32
}

func New(limit int) *Chan {
    return &Chan {
        list : make(chan interface{}, limit),
    }
}

// 将数据压入队列
func (q *Chan) Push(v interface{}) error {
    if atomic.LoadInt32(&q.closed) > 0 {
        return errors.New("closed")
    }
    q.list <- v
    return nil
}

// 先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
// 第二个返回值表示队列是否关闭
func (q *Chan) Pop() interface{} {
    return <- q.list
}

// 关闭队列(通知所有通过Pop阻塞的协程退出)
func (q *Chan) Close() {
    if atomic.LoadInt32(&q.closed) == 0 {
        atomic.StoreInt32(&q.closed, 1)
        close(q.list)
    }
}

// 获取当前队列大小
func (q *Chan) Size() int {
    return len(q.list)
}
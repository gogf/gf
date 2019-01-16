// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gchan provides graceful operations for channel.
//
// 优雅的Channel操作.
package gchan

import (
    "errors"
    "gitee.com/johng/gf/g/container/gtype"
)

type Chan struct {
    list   chan interface{}
    closed *gtype.Bool
}

func New(limit int) *Chan {
    return &Chan {
        list   : make(chan interface{}, limit),
        closed : gtype.NewBool(),
    }
}

// 将数据压入队列
func (q *Chan) Push(v interface{}) error {
    if q.closed.Val() {
        return errors.New("closed")
    }
    q.list <- v
    return nil
}

// 先进先出地从队列取出一项数据，当没有数据可获取时，阻塞等待
func (q *Chan) Pop() interface{} {
    return <- q.list
}

// 关闭队列(通知所有通过Pop阻塞的协程退出)
func (q *Chan) Close() {
    if !q.closed.Set(true) {
        close(q.list)
    }
}

// 获取当前队列大小
func (q *Chan) Size() int {
    return len(q.list)
}
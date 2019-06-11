<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 并发安全的动态队列.
// 优点：
// 1、队列初始化速度快；
// 2、可以向队头/队尾进行Push/Pop操作；
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
=======
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gqueue provides a dynamic/static concurrent-safe queue.
//
// Features:
//
// 1. FIFO queue(data -> list -> chan);
//
// 2. Fast creation and initialization;
//
// 3. Support dynamic queue size(unlimited queue size);
//
// 4. Blocking when reading data from queue;
//
package gqueue

import (
    "github.com/gogf/gf/g/container/glist"
    "math"
)

type Queue struct {
    limit     int              // Limit for queue size.
    list      *glist.List      // Underlying list structure for data maintaining.
    events    chan struct{}    // Events for data writing.
    closed    chan struct{}    // Events for queue closing.
    C         chan interface{} // Underlying channel for data reading.
}

const (
    // Size for queue buffer.
    gDEFAULT_QUEUE_SIZE = 10000
)

// New returns an empty queue object.
// Optional parameter <limit> is used to limit the size of the queue, which is unlimited in default.
// When <limit> is given, the queue will be static and high performance which is comparable with stdlib channel.
func New(limit...int) *Queue {
    q := &Queue {
        closed : make(chan struct{}, 0),
    }
    if len(limit) > 0 {
        q.limit  = limit[0]
        q.C      = make(chan interface{}, limit[0])
    } else {
        q.list   = glist.New()
        q.events = make(chan struct{}, math.MaxInt32)
        q.C      = make(chan interface{}, gDEFAULT_QUEUE_SIZE)
        go q.startAsyncLoop()
    }
    return q
}

// startAsyncLoop starts an asynchronous goroutine,
// which handles the data synchronization from list <q.list> to channel <q.C>.
func (q *Queue) startAsyncLoop() {
    for {
        select {
            case <- q.closed:
                return
            case <- q.events:
                for {
                    if length := q.list.Len(); length > 0 {
                        array := make([]interface{}, length)
                        for i := 0; i < length; i++ {
                            if e := q.list.Front(); e != nil {
                                array[i] = q.list.Remove(e)
                            } else {
                                break
                            }
                        }
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

// Push pushes the data <v> into the queue.
// Note that it would panics if Push is called after the queue is closed.
func (q *Queue) Push(v interface{}) {
    if q.limit > 0 {
        q.C <- v
    } else {
        q.list.PushBack(v)
        q.events <- struct{}{}
    }
}

// Pop pops an item from the queue in FIFO way.
// Note that it would return nil immediately if Pop is called after the queue is closed.
func (q *Queue) Pop() interface{} {
    return <- q.C
}

// Close closes the queue.
// Notice: It would notify all goroutines return immediately,
// which are being blocked reading using Pop method.
func (q *Queue) Close() {
    close(q.C)
    close(q.events)
    close(q.closed)
}

// Size returns the length of the queue.
func (q *Queue) Size() int {
    return len(q.C) + q.list.Len()
>>>>>>> upstream/master
}



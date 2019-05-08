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
// Optional parameter <limit> is used to limit the size of the queue, which is unlimited by default.
// When <limit> is given, the queue will be static and high performance which is comparable with stdlib chan.
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
// which are being blocked reading by Pop method.
func (q *Queue) Close() {
    close(q.C)
    close(q.events)
    close(q.closed)
}

// Size returns the length of the queue.
func (q *Queue) Size() int {
    return len(q.C) + q.list.Len()
}



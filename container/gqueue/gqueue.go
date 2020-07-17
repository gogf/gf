// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gqueue provides dynamic/static concurrent-safe queue.
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
	"math"

	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/container/gtype"
)

// Queue is a concurrent-safe queue built on doubly linked list and channel.
type Queue struct {
	limit  int              // Limit for queue size.
	list   *glist.List      // Underlying list structure for data maintaining.
	closed *gtype.Bool      // Whether queue is closed.
	events chan struct{}    // Events for data writing.
	C      chan interface{} // Underlying channel for data reading.
}

const (
	// Size for queue buffer.
	gDEFAULT_QUEUE_SIZE = 10000
	// Max batch size per-fetching from list.
	gDEFAULT_MAX_BATCH_SIZE = 10
)

// New returns an empty queue object.
// Optional parameter <limit> is used to limit the size of the queue, which is unlimited in default.
// When <limit> is given, the queue will be static and high performance which is comparable with stdlib channel.
func New(limit ...int) *Queue {
	q := &Queue{
		closed: gtype.NewBool(),
	}
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan interface{}, limit[0])
	} else {
		q.list = glist.New(true)
		q.events = make(chan struct{}, math.MaxInt32)
		q.C = make(chan interface{}, gDEFAULT_QUEUE_SIZE)
		go q.asyncLoopFromListToChannel()
	}
	return q
}

// asyncLoopFromListToChannel starts an asynchronous goroutine,
// which handles the data synchronization from list <q.list> to channel <q.C>.
func (q *Queue) asyncLoopFromListToChannel() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.events
		for !q.closed.Val() {
			if length := q.list.Len(); length > 0 {
				if length > gDEFAULT_MAX_BATCH_SIZE {
					length = gDEFAULT_MAX_BATCH_SIZE
				}
				for _, v := range q.list.PopFronts(length) {
					// When q.C is closed, it will panic here, especially q.C is being blocked for writing.
					// If any error occurs here, it will be caught by recover and be ignored.
					q.C <- v
				}
			} else {
				break
			}
		}
		// Clear q.events to remain just one event to do the next synchronization check.
		for i := 0; i < len(q.events)-1; i++ {
			<-q.events
		}
	}
	// It should be here to close q.C if <q> is unlimited size.
	// It's the sender's responsibility to close channel when it should be closed.
	close(q.C)
}

// Push pushes the data <v> into the queue.
// Note that it would panics if Push is called after the queue is closed.
func (q *Queue) Push(v interface{}) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		if len(q.events) < gDEFAULT_QUEUE_SIZE {
			q.events <- struct{}{}
		}
	}
}

// Pop pops an item from the queue in FIFO way.
// Note that it would return nil immediately if Pop is called after the queue is closed.
func (q *Queue) Pop() interface{} {
	return <-q.C
}

// Close closes the queue.
// Notice: It would notify all goroutines return immediately,
// which are being blocked reading using Pop method.
func (q *Queue) Close() {
	q.closed.Set(true)
	if q.events != nil {
		close(q.events)
	}
	if q.limit > 0 {
		close(q.C)
	}
	for i := 0; i < gDEFAULT_MAX_BATCH_SIZE; i++ {
		q.Pop()
	}
}

// Len returns the length of the queue.
// Note that the result might not be accurate as there's a
// asynchronize channel reading the list constantly.
func (q *Queue) Len() (length int) {
	if q.list != nil {
		length += q.list.Len()
	}
	length += len(q.C)
	return
}

// Size is alias of Len.
func (q *Queue) Size() int {
	return q.Len()
}

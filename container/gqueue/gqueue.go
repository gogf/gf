// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
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
package gqueue

import (
	"math"

	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/container/gtype"
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
	defaultQueueSize = 10000 // Size for queue buffer.
	defaultBatchSize = 10    // Max batch size per-fetching from list.
)

// New returns an empty queue object.
// Optional parameter `limit` is used to limit the size of the queue, which is unlimited in default.
// When `limit` is given, the queue will be static and high performance which is comparable with stdlib channel.
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
		q.C = make(chan interface{}, defaultQueueSize)
		go q.asyncLoopFromListToChannel()
	}
	return q
}

// Push pushes the data `v` into the queue.
// Note that it would panic if Push is called after the queue is closed.
func (q *Queue) Push(v interface{}) {
	if q.limit > 0 {
		q.C <- v
	} else {
		q.list.PushBack(v)
		if len(q.events) < defaultQueueSize {
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
	if !q.closed.Cas(false, true) {
		return
	}
	if q.events != nil {
		close(q.events)
	}
	if q.limit > 0 {
		close(q.C)
	} else {
		for i := 0; i < defaultBatchSize; i++ {
			q.Pop()
		}
	}
}

// Len returns the length of the queue.
// Note that the result might not be accurate if using unlimited queue size as there's an
// asynchronous channel reading the list constantly.
func (q *Queue) Len() (length int64) {
	bufferedSize := int64(len(q.C))
	if q.limit > 0 {
		return bufferedSize
	}
	return int64(q.list.Size()) + bufferedSize
}

// Size is alias of Len.
// Deprecated, use Len instead.
func (q *Queue) Size() int64 {
	return q.Len()
}

// asyncLoopFromListToChannel starts an asynchronous goroutine,
// which handles the data synchronization from list `q.list` to channel `q.C`.
func (q *Queue) asyncLoopFromListToChannel() {
	defer func() {
		if q.closed.Val() {
			_ = recover()
		}
	}()
	for !q.closed.Val() {
		<-q.events
		for !q.closed.Val() {
			if bufferLength := q.list.Len(); bufferLength > 0 {
				// When q.C is closed, it will panic here, especially q.C is being blocked for writing.
				// If any error occurs here, it will be caught by recover and be ignored.
				for i := 0; i < bufferLength; i++ {
					q.C <- q.list.PopFront()
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
	// It should be here to close `q.C` if `q` is unlimited size.
	// It's the sender's responsibility to close channel when it should be closed.
	close(q.C)
}

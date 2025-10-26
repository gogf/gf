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

// Queue is a concurrent-safe queue built on doubly linked list and channel.
type Queue struct {
	*TQueue[any]
}

const (
	defaultQueueSize = 10000 // Size for queue buffer.
	defaultBatchSize = 10    // Max batch size per-fetching from list.
)

// New returns an empty queue object.
// Optional parameter `limit` is used to limit the size of the queue, which is unlimited in default.
// When `limit` is given, the queue will be static and high performance which is comparable with stdlib channel.
func New(limit ...int) *Queue {
	return &Queue{
		TQueue: NewTQueue[any](limit...),
	}
}

// Push pushes the data `v` into the queue.
// Note that it would panic if Push is called after the queue is closed.
func (q *Queue) Push(v any) {
	q.TQueue.Push(v)
}

// Pop pops an item from the queue in FIFO way.
// Note that it would return nil immediately if Pop is called after the queue is closed.
func (q *Queue) Pop() any {
	return q.TQueue.Pop()
}

// Close closes the queue.
// Notice: It would notify all goroutines return immediately,
// which are being blocked reading using Pop method.
func (q *Queue) Close() {
	q.TQueue.Close()
}

// Len returns the length of the queue.
// Note that the result might not be accurate if using unlimited queue size as there's an
// asynchronous channel reading the list constantly.
func (q *Queue) Len() (length int64) {
	return q.TQueue.Len()
}

// Size is alias of Len.
// Deprecated: use Len instead.
func (q *Queue) Size() int64 {
	return q.Len()
}

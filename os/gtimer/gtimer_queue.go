// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

import (
	"container/heap"
	"math"
	"sync"

	"github.com/gogf/gf/v3/container/gatomic"
)

// priorityQueue is an abstract data type similar to a regular queue or stack data structure in which
// each element additionally has a "priority" associated with it. In a priority queue, an element with
// high priority is served before an element with low priority.
// priorityQueue is based on heap structure.
type priorityQueue struct {
	mu           sync.Mutex
	heap         *priorityQueueHeap // the underlying queue items manager using heap.
	nextPriority *gatomic.Int64     // nextPriority stores the next priority value of the heap, which is used to check if necessary to call the Pop of heap by Timer.
}

// priorityQueueHeap is a heap manager, of which the underlying `array` is an array implementing a heap structure.
type priorityQueueHeap struct {
	array []priorityQueueItem
}

// priorityQueueItem stores the queue item which has a `priority` attribute to sort itself in heap.
type priorityQueueItem struct {
	value    interface{}
	priority int64
}

// newPriorityQueue creates and returns a priority queue.
func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{
		heap:         &priorityQueueHeap{array: make([]priorityQueueItem, 0)},
		nextPriority: gatomic.NewInt64(math.MaxInt64),
	}
	heap.Init(queue.heap)
	return queue
}

// NextPriority retrieves and returns the minimum and the most priority value of the queue.
func (q *priorityQueue) NextPriority() int64 {
	return q.nextPriority.Val()
}

// Push pushes a value to the queue.
// The `priority` specifies the priority of the value.
// The lesser the `priority` value the higher priority of the `value`.
func (q *priorityQueue) Push(value interface{}, priority int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	heap.Push(q.heap, priorityQueueItem{
		value:    value,
		priority: priority,
	})
	// Update the minimum priority using atomic operation.
	nextPriority := q.nextPriority.Val()
	if priority >= nextPriority {
		return
	}
	q.nextPriority.Set(priority)
}

// Pop retrieves, removes and returns the most high priority value from the queue.
func (q *priorityQueue) Pop() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if v := heap.Pop(q.heap); v != nil {
		var nextPriority int64 = math.MaxInt64
		if len(q.heap.array) > 0 {
			nextPriority = q.heap.array[0].priority
		}
		q.nextPriority.Set(nextPriority)
		return v.(priorityQueueItem).value
	}
	return nil
}

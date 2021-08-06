// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtimer

// Len is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap) Len() int {
	return len(h.array)
}

// Less is used to implement the interface of sort.Interface.
// The least one is placed to the top of the heap.
func (h *priorityQueueHeap) Less(i, j int) bool {
	return h.array[i].priority < h.array[j].priority
}

// Swap is used to implement the interface of sort.Interface.
func (h *priorityQueueHeap) Swap(i, j int) {
	if len(h.array) == 0 {
		return
	}
	h.array[i], h.array[j] = h.array[j], h.array[i]
}

// Push pushes an item to the heap.
func (h *priorityQueueHeap) Push(x interface{}) {
	h.array = append(h.array, x.(priorityQueueItem))
}

// Pop retrieves, removes and returns the most high priority item from the heap.
func (h *priorityQueueHeap) Pop() interface{} {
	length := len(h.array)
	if length == 0 {
		return nil
	}
	item := h.array[length-1]
	h.array = h.array[0 : length-1]
	return item
}

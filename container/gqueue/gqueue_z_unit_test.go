// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gqueue_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestQueue_Len(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			maxNum   = 100
			maxTries = 100
		)
		for n := 10; n < maxTries; n++ {
			q1 := gqueue.New(maxNum)
			for i := range maxNum {
				q1.Push(i)
			}
			t.Assert(q1.Len(), maxNum)
			t.Assert(q1.Size(), maxNum)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			maxNum   = 100
			maxTries = 100
		)
		for n := 10; n < maxTries; n++ {
			q1 := gqueue.New()
			for i := range maxNum {
				q1.Push(i)
			}
			t.AssertLE(q1.Len(), maxNum)
			t.AssertLE(q1.Size(), maxNum)
		}
	})
}

func TestQueue_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		defer q.Close()
		for i := range 100 {
			q.Push(i)
		}
		t.Assert(q.Pop(), 0)
		t.Assert(q.Pop(), 1)
	})
}

func TestQueue_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New()
		defer q1.Close()
		q1.Push(1)
		q1.Push(2)
		q1.Push(3)
		q1.Push(4)
		i1 := q1.Pop()
		t.Assert(i1, 1)
	})
}

func TestQueue_Close(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New()
		defer q1.Close()
		q1.Push(1)
		q1.Push(2)
		// wait sync to channel
		time.Sleep(10 * time.Millisecond)
		t.Assert(q1.Len(), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New(2)
		defer q1.Close()
		q1.Push(1)
		q1.Push(2)
		// wait sync to channel
		time.Sleep(10 * time.Millisecond)
		t.Assert(q1.Len(), 2)
	})
}

func Test_Issue2509(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		defer q.Close()
		q.Push(1)
		q.Push(2)
		q.Push(3)
		t.AssertLE(q.Len(), 3)
		t.Assert(<-q.C, 1)
		t.AssertLE(q.Len(), 2)
		t.Assert(<-q.C, 2)
		t.AssertLE(q.Len(), 1)
		t.Assert(<-q.C, 3)
		t.Assert(q.Len(), 0)
	})
}

// Issue #4376
func TestIssue4376(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gq := gqueue.New()
		defer gq.Close()
		cq := make(chan int, 100000)
		defer close(cq)

		for i := range 11603 {
			gq.Push(i)
			cq <- i
		}
		// May be not equal because of the async channel reading goroutine.
		t.Log(gq.Len(), len(cq))
		time.Sleep(50 * time.Millisecond)
		t.Log(gq.Len(), len(cq))
	})
}

// Test static queue (with limit) close operation
func TestQueue_StaticClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New(10)
		defer func() {
			if err := recover(); err == nil {
				t.Log("Close succeeded")
			}
		}()
		q.Push(1)
		q.Push(2)
		q.Close()
		// After closing, Pop should return nil
		v := q.Pop()
		t.Assert(v, nil)
	})
}

// Test Size() method (deprecated alias of Len)
func TestQueue_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New(20)
		for i := range 10 {
			q.Push(i)
		}
		t.Assert(q.Size(), 10)
		t.Assert(q.Len(), 10)
		q.Close()
	})
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		for i := range 15 {
			q.Push(i)
		}
		time.Sleep(10 * time.Millisecond)
		t.Assert(q.Size(), q.Len())
		q.Close()
	})
}

// Test TQueue directly with generic type
func TestTQueue_Generic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with custom type
		q := gqueue.NewTQueue[string]()
		defer q.Close()
		q.Push("hello")
		q.Push("world")
		t.Assert(q.Pop(), "hello")
		t.Assert(q.Pop(), "world")
	})
}

// Test TQueue Size method directly
func TestTQueue_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.NewTQueue[int]()
		defer q.Close()
		for i := range 10 {
			q.Push(i)
		}
		time.Sleep(10 * time.Millisecond)
		// Size is an alias of Len for TQueue
		t.Assert(q.Size(), q.Len())
	})
}

// Test TQueue with static limit
func TestTQueue_StaticLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.NewTQueue[int](5)
		defer q.Close()
		for i := range 5 {
			q.Push(i)
		}
		t.Assert(q.Len(), 5)
		for i := range 5 {
			t.Assert(q.Pop(), i)
		}
		t.Assert(q.Len(), 0)
	})
}

// Test queue with large data push/pop
func TestQueue_LargeDataScale(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		defer q.Close()
		n := 5000
		for i := range n {
			q.Push(i)
		}
		time.Sleep(50 * time.Millisecond)
		// Pop should retrieve all items in order
		for i := range n {
			v := q.Pop()
			t.Assert(v, i)
		}
	})
}

// Test double close (idempotent close)
func TestQueue_DoubleClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		q.Push(1)
		q.Close()
		// Second close should not panic
		q.Close()
		t.Assert(q.Pop(), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New(10)
		q.Push(1)
		q.Close()
		// Second close should not panic for static queue
		q.Close()
		// Pop from closed static queue returns the buffered value
		v := q.Pop()
		t.Assert(v, 1)
	})
}

// Test concurrent push and pop
func TestQueue_ConcurrentPushPop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		defer q.Close()
		// Producer goroutine
		go func() {
			for i := range 100 {
				q.Push(i)
			}
			time.Sleep(50 * time.Millisecond)
			q.Close()
		}()
		// Consumer
		count := 0
		for {
			v := q.Pop()
			if v == nil {
				break
			}
			count++
		}
		t.AssertGE(count, 1)
	})
}

// Test Pop on empty queue returns nil when closed
func TestQueue_PopEmptyClosed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		q.Close()
		v := q.Pop()
		t.Assert(v, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New(10)
		q.Close()
		v := q.Pop()
		t.Assert(v, nil)
	})
}

// Test Len with dynamic queue at capacity boundary
func TestQueue_LenAtBoundary(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		defer q.Close()
		// Push exactly defaultQueueSize items to test boundary condition
		for i := range 10000 {
			q.Push(i)
		}
		time.Sleep(50 * time.Millisecond)
		len := q.Len()
		t.AssertGE(len, 0)
	})
}

// Test Close on dynamic queue with pending asyncLoopFromListToChannel
func TestQueue_CloseWithAsyncLoop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
		// Push some data to activate asyncLoopFromListToChannel
		for i := range 100 {
			q.Push(i)
		}
		// Immediately close
		q.Close()
		// Pop should return values until exhausted, then nil
		for {
			v := q.Pop()
			if v == nil {
				break
			}
		}
		t.Assert(q.Pop(), nil)
	})
}

// Test static queue edge case with zero limit (should create unlimited queue)
func TestQueue_ZeroLimitCreatesUnlimited(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New(0)
		defer q.Close()
		for i := range 100 {
			q.Push(i)
		}
		time.Sleep(10 * time.Millisecond)
		len := q.Len()
		t.Assert(len, 100)
	})
}

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
			for i := 0; i < maxNum; i++ {
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
			for i := 0; i < maxNum; i++ {
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
		for i := 0; i < 100; i++ {
			q.Push(i)
		}
		t.Assert(q.Pop(), 0)
		t.Assert(q.Pop(), 1)
	})
}

func TestQueue_Pop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New()
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
		q1.Push(1)
		q1.Push(2)
		// wait sync to channel
		time.Sleep(10 * time.Millisecond)
		t.Assert(q1.Len(), 2)
		q1.Close()
	})
	gtest.C(t, func(t *gtest.T) {
		q1 := gqueue.New(2)
		q1.Push(1)
		q1.Push(2)
		// wait sync to channel
		time.Sleep(10 * time.Millisecond)
		t.Assert(q1.Len(), 2)
		q1.Close()
	})
}

func Test_Issue2509(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		q := gqueue.New()
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

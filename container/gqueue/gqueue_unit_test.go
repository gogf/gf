// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gqueue_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/test/gtest"
)

func TestQueue_Len(t *testing.T) {
	max := 100
	for n := 10; n < max; n++ {
		q1 := gqueue.New(max)
		for i := 0; i < max; i++ {
			q1.Push(i)
		}
		gtest.Assert(q1.Len(), max)
		gtest.Assert(q1.Size(), max)
	}
}

func TestQueue_Basic(t *testing.T) {
	q := gqueue.New()
	for i := 0; i < 100; i++ {
		q.Push(i)
	}
	gtest.Assert(q.Pop(), 0)
	gtest.Assert(q.Pop(), 1)
}

func TestQueue_Pop(t *testing.T) {
	q1 := gqueue.New()
	q1.Push(1)
	q1.Push(2)
	q1.Push(3)
	q1.Push(4)
	i1 := q1.Pop()
	gtest.Assert(i1, 1)
}

func TestQueue_Close(t *testing.T) {
	q1 := gqueue.New()
	q1.Push(1)
	q1.Push(2)
	time.Sleep(time.Millisecond)
	gtest.Assert(q1.Len(), 2)
	q1.Close()
}

package gqueue_test

import (
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestQueue_Len(t *testing.T) {
	q1 := gqueue.New(2)
	q1.Push(1)
	q1.Push(2)
	gtest.Assert(q1.Len(), 2)
	gtest.Assert(q1.Size(), 2)
}

func TestQueue_Pop(t *testing.T) {
	q1 := gqueue.New()
	q1.Push(1)
	q1.Push(2)
	q1.Push(3)
	i1 := q1.Pop()
	gtest.Assert(i1, 1)
}

func TestQueue_Close(t *testing.T) {
	q1 := gqueue.New()
	q1.Push(1)
	q1.Push(2)
	gtest.Assert(q1.Len(), 2)
	q1.Close()
	gtest.Assert(q1.Len(), 2)
}

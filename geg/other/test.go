package main

import (
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	max := 100
	q := gqueue.New(max)
	for i := 1; i < max; i++ {
		q.Push(i)
	}
	q.Close()
	gtest.Assert(q.Len(), 1)

}

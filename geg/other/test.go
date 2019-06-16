package main

import (
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	max := 100
	for i := 0; i < 10; i++ {
		q := gqueue.New(max)
		for i := 1; i < max; i++ {
			q.Push(i)
		}
		gtest.Assert(q.Pop(), 1)
	}
}

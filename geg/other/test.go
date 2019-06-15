package main

import (
	"github.com/gogf/gf/g/container/gqueue"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	q1 := gqueue.New(2)
	q1.Push(1)
	q1.Push(2)
	gtest.Assert(q1.Size(),2)
}


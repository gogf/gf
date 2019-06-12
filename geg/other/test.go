package main

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	a2 := []int{1}
	array2 := garray.NewSortedIntArrayFrom(a2)
	gtest.Assert(array2.Search(2),-1)
}
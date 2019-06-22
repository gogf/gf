package main

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	a2 := []interface{}{0, 1, 2, 3, 4, 5, 6}
	array3 := garray.NewArrayFrom(a2, true)
	gtest.Assert(array3.SubSlice(2, 2), []interface{}{2, 3})
}

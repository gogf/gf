package main

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {
	a1:=[]string{"a", "d", "c","b"}

	s1 :=garray.NewSortedStringArrayFromCopy(a1,true)

	gtest.Assert(s1.Slice(),[]string{"a", "b", "c","d"})
}

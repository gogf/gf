package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gtype"
	"github.com/gogf/gf/g/test/gtest"
)

func main() {

	var add float32
	add = 3.1415926

	myF32 := gtype.NewFloat32(add)
	myAdd := myF32.Add(float32(6.2951413))

	add += float32(6.2951413)

	fmt.Println(myF32.Val())
	gtest.AssertEQ(myAdd, add)

}
package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

func main() {
	var v g.Var

	v.Set("123")

	fmt.Println(v.Val())

	// 基本类型转换
	fmt.Println(v.Int())
	fmt.Println(v.Uint())
	fmt.Println(v.Float64())

	// slice转换
	fmt.Println(v.Ints())
	fmt.Println(v.Floats())
	fmt.Println(v.Strings())

	// struct转换
	type Score struct {
		Value int
	}
	s := new(Score)
	v.Struct(s)
	fmt.Println(s)
}

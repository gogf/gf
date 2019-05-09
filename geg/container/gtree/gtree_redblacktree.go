package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gtree"
)

func main() {
	tree := gtree.New(func(v1, v2 interface{}) int {
		return v1.(int) - v2.(int)
	})
	for i := 0; i < 20; i++ {
		tree.Set(i, i)
	}
	tree.Print()
	tree.IteratorAsc(func(key, value interface{}) bool {
		fmt.Println(key)
		return true
	})
	fmt.Println()
	tree.IteratorDesc(func(key, value interface{}) bool {
		fmt.Println(key)
		return true
	})
}
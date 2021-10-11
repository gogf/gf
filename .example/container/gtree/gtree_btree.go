package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gtree"
)

func main() {
	tree := gtree.NewBTree(10, func(v1, v2 interface{}) int {
		return v1.(int) - v2.(int)
	})
	for i := 0; i < 20; i++ {
		tree.Set(i, i*10)
	}
	fmt.Println(tree.String())

	tree.IteratorDesc(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

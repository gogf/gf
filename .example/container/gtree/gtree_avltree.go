package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/util/gutil"
)

func main() {
	tree := gtree.NewAVLTree(gutil.ComparatorInt)
	for i := 0; i < 10; i++ {
		tree.Set(i, i*10)
	}
	// 打印树形
	tree.Print()
	// 前序遍历
	fmt.Println("ASC:")
	tree.IteratorAsc(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
	// 后续遍历
	fmt.Println("DESC:")
	tree.IteratorDesc(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}

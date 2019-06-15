package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/gtree"
	"github.com/gogf/gf/g/util/gutil"
)

func main() {
	expect := map[interface{}]interface{}{
		20: "val20",
		6:  "val6",
		10: "val10",
		12: "val12",
		1:  "val1",
		15: "val15",
		19: "val19",
		8:  "val8",
		4:  "val4"}
	m := gtree.NewAVLTreeFrom(gutil.ComparatorInt, expect)
	m.Print()

	//m := avltree.NewWithIntComparator()
	//m.Remove()
	fmt.Println(1, m.Remove(1))// 应该输出val1，但输出nil
	fmt.Println(2, m.Remove(1))
	fmt.Println(3, m.Get(1))

	fmt.Println(4, m.Remove(20))// 应该输出val20，但输出nil
	fmt.Println(5, m.Remove(20))
	fmt.Println(6, m.Get(20))
}


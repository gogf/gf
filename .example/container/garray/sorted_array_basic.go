package main

import (
	"fmt"

	"github.com/jin502437344/gf/container/garray"
)

func main() {
	// 自定义排序数组，降序排序(SortedIntArray管理的数据是升序)
	a := garray.NewSortedArray(func(v1, v2 interface{}) int {
		if v1.(int) < v2.(int) {
			return 1
		}
		if v1.(int) > v2.(int) {
			return -1
		}
		return 0
	})

	// 添加数据
	a.Add(2)
	a.Add(3)
	a.Add(1)
	fmt.Println(a.Slice())

	// 添加重复数据
	a.Add(3)
	fmt.Println(a.Slice())

	// 检索数据，返回最后对比的索引位置，检索结果
	// 检索结果：0: 匹配; <0:参数小于对比值; >0:参数大于对比值
	fmt.Println(a.Search(1))

	// 设置不可重复
	a.SetUnique(true)
	fmt.Println(a.Slice())
	a.Add(1)
	fmt.Println(a.Slice())
}

package main

import (
	"fmt"
	"github.com/gogf/gf/g/container/garray"
)

func main() {
	// 创建普通的int类型数组，并关闭默认的并发安全特性
	a := garray.NewIntArray(true)

	// 添加数据项
	for i := 0; i < 10; i++ {
		a.Append(i)
	}

	// 获取当前数组长度
	fmt.Println(a.Len())

	// 获取当前数据项列表
	fmt.Println(a.Slice())

	// 获取指定索引项
	fmt.Println(a.Get(6))

	// 在指定索引前插入数据项
	a.InsertAfter(9, 11)
	// 在指定索引后插入数据项
	a.InsertBefore(10, 10)
	fmt.Println(a.Slice())

	// 修改指定索引的数据项
	a.Set(0, 100)
	fmt.Println(a.Slice())

	// 搜索数据项，返回搜索到的索引位置
	fmt.Println(a.Search(5))

	// 删除指定索引的数据项
	a.Remove(0)
	fmt.Println(a.Slice())

	// 并发安全，写锁操作
	a.LockFunc(func(array []int) {
		// 将末尾项改为100
		array[len(array)-1] = 100
	})

	// 并发安全，读锁操作
	a.RLockFunc(func(array []int) {
		fmt.Println(array[len(array)-1])
	})

	// 清空数组
	fmt.Println(a.Slice())
	a.Clear()
	fmt.Println(a.Slice())
}

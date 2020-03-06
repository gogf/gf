// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"

	"github.com/gogf/gf/container/garray"
)

func Example_basic() {
	// 创建普通的数组
	a := garray.New()

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

	// 查找指定数据项是否存在
	fmt.Println(a.Contains(6))
	fmt.Println(a.Contains(100))

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

	// 清空数组
	fmt.Println(a.Slice())
	a.Clear()
	fmt.Println(a.Slice())

	// Output:
	// 10
	// [0 1 2 3 4 5 6 7 8 9]
	// 6
	// true
	// false
	// [0 1 2 3 4 5 6 7 8 9 10 11]
	// [100 1 2 3 4 5 6 7 8 9 10 11]
	// 5
	// [1 2 3 4 5 6 7 8 9 10 11]
	// [1 2 3 4 5 6 7 8 9 10 11]
	// []
}

func Example_rand() {
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})
	// 随机返回两个数据项(不删除)
	fmt.Println(array.Rands(2))
	fmt.Println(array.PopRand())
}

func Example_pop() {
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})
	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1
	// [2 3]
	// 9
	// [7 8]
}

func Example_merge() {
	array1 := garray.NewFrom([]interface{}{1, 2})
	array2 := garray.NewFrom([]interface{}{3, 4})
	slice1 := []interface{}{5, 6}
	slice2 := []int{7, 8}
	slice3 := []string{"9", "0"}
	fmt.Println(array1.Slice())
	array1.Merge(array1)
	array1.Merge(array2)
	array1.Merge(slice1)
	array1.Merge(slice2)
	array1.Merge(slice3)
	fmt.Println(array1.Slice())

	// Output:
	// [1 2]
	// [1 2 1 2 3 4 5 6 7 8 9 0]
}

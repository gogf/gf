package main

import (
	"fmt"

	"github.com/jin502437344/gf/container/gtree"
	"github.com/jin502437344/gf/util/gutil"
)

func main() {
	m := gtree.NewRedBlackTree(gutil.ComparatorInt)

	// 设置键值对
	for i := 0; i < 10; i++ {
		m.Set(i, i*10)
	}
	// 查询大小
	fmt.Println(m.Size())
	// 批量设置键值对(不同的数据类型对象参数不同)
	m.Sets(map[interface{}]interface{}{
		10: 10,
		11: 11,
	})
	fmt.Println(m.Size())

	// 查询是否存在
	fmt.Println(m.Contains(1))

	// 查询键值
	fmt.Println(m.Get(1))

	// 删除数据项
	m.Remove(9)
	fmt.Println(m.Size())

	// 批量删除
	m.Removes([]interface{}{10, 11})
	fmt.Println(m.Size())

	// 当前键名列表(随机排序)
	fmt.Println(m.Keys())
	// 当前键值列表(随机排序)
	fmt.Println(m.Values())

	// 查询键名，当键值不存在时，写入给定的默认值
	fmt.Println(m.GetOrSet(100, 100))

	// 删除键值对，并返回对应的键值
	fmt.Println(m.Remove(100))

	// 遍历map
	m.IteratorAsc(func(k interface{}, v interface{}) bool {
		fmt.Printf("%v:%v ", k, v)
		return true
	})
	fmt.Println()

	// 清空map
	m.Clear()

	// 判断map是否为空
	fmt.Println(m.IsEmpty())
}

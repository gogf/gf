package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gmap"
)

func main() {
	// 创建一个默认的gmap对象，
	// 默认情况下该gmap对象不支持并发安全特性，
	// 初始化时可以给定true参数开启并发安全特性，用以并发安全场景。
	m := gmap.New()

	// 设置键值对
	for i := 0; i < 10; i++ {
		m.Set(i, i)
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
	m.Iterator(func(k interface{}, v interface{}) bool {
		fmt.Printf("%v:%v ", k, v)
		return true
	})

	// 自定义写锁操作
	m.LockFunc(func(m map[interface{}]interface{}) {
		m[99] = 99
	})

	// 自定义读锁操作
	m.RLockFunc(func(m map[interface{}]interface{}) {
		fmt.Println(m[99])
	})

	// 清空map
	m.Clear()

	// 判断map是否为空
	fmt.Println(m.IsEmpty())
}

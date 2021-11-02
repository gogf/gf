package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
)

func main() {
	// 创建一个非并发安全的集合对象
	s := gset.New(true)

	// 添加数据项
	s.Add(1)

	// 批量添加数据项
	s.Add([]interface{}{1, 2, 3}...)

	// 集合数据项大小
	fmt.Println(s.Size())

	// 集合中是否存在指定数据项
	fmt.Println(s.Contains(2))

	// 返回数据项slice
	fmt.Println(s.Slice())

	// 删除数据项
	s.Remove(3)

	// 遍历数据项
	s.Iterator(func(v interface{}) bool {
		fmt.Println("Iterator:", v)
		return true
	})

	// 将集合转换为字符串
	fmt.Println(s.String())

	// 并发安全写锁操作
	s.LockFunc(func(m map[interface{}]struct{}) {
		m[4] = struct{}{}
	})

	// 并发安全读锁操作
	s.RLockFunc(func(m map[interface{}]struct{}) {
		fmt.Println(m)
	})

	// 清空集合
	s.Clear()
	fmt.Println(s.Size())
}

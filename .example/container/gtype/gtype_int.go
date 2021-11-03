package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gtype"
)

func main() {
	// 创建一个Int型的并发安全基本类型对象
	i := gtype.NewInt()

	// 设置值
	i.Set(10)

	// 获取值
	fmt.Println(i.Val())

	// (整型/浮点型有效)数值 增加/删除 delta
	fmt.Println(i.Add(-1))
}

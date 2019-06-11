package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
)

// 演示slice类型属性的赋值
func main() {
	type User struct {
		Scores []int
	}

	user := new(User)
	scores := []interface{}{99, 100, 60, 140}

	// 通过map映射转换
	if err := gconv.Struct(g.Map{"Scores": scores}, user); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user)
	}

	// 通过变量映射转换，直接slice赋值
	if err := gconv.Struct(scores, user); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user)
	}
}

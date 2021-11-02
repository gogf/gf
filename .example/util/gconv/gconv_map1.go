package main

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	type User struct {
		Uid  int    `json:"uid"`
		Name string `json:"name"`
	}
	// 对象
	fmt.Println(gconv.Map(User{
		Uid:  1,
		Name: "john",
	}))
	// 对象指针
	fmt.Println(gconv.Map(&User{
		Uid:  1,
		Name: "john",
	}))

	// 任意map类型
	fmt.Println(gconv.Map(map[int]int{
		100: 10000,
	}))
}

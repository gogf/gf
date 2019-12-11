package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	type User struct {
		Scores []int64
	}
	user := new(User)
	err := gconv.Struct(g.Map{"scores": []interface{}{1, 2, 3}}, user)
	fmt.Println(err)
	fmt.Println(user)
}

package main

import (
	"fmt"

	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	type User struct {
		Uid      int
		Name     string `gconv:"-"`
		NickName string `gconv:"nickname, omitempty"`
		Pass1    string `gconv:"password1"`
		Pass2    string `gconv:"password2"`
	}
	user := User{
		Uid:   100,
		Name:  "john",
		Pass1: "123",
		Pass2: "456",
	}
	fmt.Println(gconv.Map(user))
}

package main

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	type Base struct {
		Id         int    `c:"id"`
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base     `c:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	user := new(User)
	user.Id = 1
	user.Nickname = "John"
	user.Passport = "johng"
	user.Password = "123456"
	user.CreateTime = "2019"
	fmt.Println(gconv.Map(user))
	fmt.Println(gconv.MapDeep(user))
}

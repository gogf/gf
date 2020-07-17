package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
)

func main() {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	user := new(User)
	user.Id = 1
	user.Uid = 100
	user.Nickname = "John"
	user.Passport = "johng"
	user.Password = "123456"
	user.CreateTime = "2019"
	g.Dump(gconv.Map(user))
	g.Dump(gconv.MapDeep(user))
}

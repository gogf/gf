package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	type Ids struct {
		Id  int `json:"id"`
		Uid int `json:"uid"`
	}
	type Base struct {
		Ids
		CreateTime string `json:"create_time"`
	}
	type User struct {
		Base
		Passport string `json:"passport"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	data := g.Map{
		"id":          1,
		"uid":         100,
		"passport":    "johng",
		"password":    "123456",
		"nickname":    "John",
		"create_time": "2019",
	}
	user := new(User)
	gconv.StructDeep(data, user)
	g.Dump(user)
}

package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
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

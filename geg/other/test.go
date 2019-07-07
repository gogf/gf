package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gvalid"
)

func main() {
	type Pass struct {
		Pass1 string `valid:"password1@required|same:password2#请输入您的密码|您两次输入的密码不一致"`
		Pass2 string `valid:"password2@required|same:password1#请再次输入您的密码|您两次输入的密码不一致"`
	}
	type User struct {
		Id   int
		Name string `valid:"name@required#请输入您的姓名"`
		Pass Pass
	}
	user := &User{
		Name: "john",
		Pass: Pass{
			Pass1: "1",
			Pass2: "2",
		},
	}
	err := gvalid.CheckStruct(user, nil)
	g.Dump(err.Maps())
}

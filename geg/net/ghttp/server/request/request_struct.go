package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type User struct {
	Uid   int    `json:"uid"`
	Name  string `json:"name"  params:"username"`
	Pass1 string `json:"pass1" params:"password1,userpass1"`
	Pass2 string `json:"pass2" params:"password3,userpass2"`
}

func main() {
	s := g.Server()
	s.BindHandler("/user", func(r *ghttp.Request) {
		user := new(User)
		r.GetToStruct(user)
		//r.GetPostToStruct(user)
		//r.GetQueryToStruct(user)
		r.Response.WriteJson(user)
	})
	s.SetPort(8199)
	s.Run()

	// http://127.0.0.1:8199/user?uid=1&name=john&password1=123&userpass2=123
	// {"name":"john","pass1":"123","pass2":"123","uid":1}
}

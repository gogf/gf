package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	type User struct {
		Uid   int    `json:"uid"`
		Name  string `json:"name"  p:"username"`
		Pass1 string `json:"pass1" p:"password1"`
		Pass2 string `json:"pass2" p:"password2"`
	}

	s := g.Server()
	s.BindHandler("/user", func(r *ghttp.Request) {
		var user *User
		if err := r.Parse(&user); err != nil {
			panic(err)
		}
		r.Response.WriteJson(user)
	})
	s.SetPort(8199)
	s.Run()

	// http://127.0.0.1:8199/user?uid=1&name=john&password1=123&userpass2=123
	// {"name":"john","pass1":"123","pass2":"123","uid":1}
}

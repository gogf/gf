package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	type User struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Pass1 string `json:"password1" p:"password1"`
		Pass2 string `json:"password2" p:"password2"`
	}
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		var user *User
		if err := r.Parse(&user); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(user)
	})
	s.SetPort(8199)
	s.Run()
}

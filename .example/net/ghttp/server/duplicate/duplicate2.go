// 路由重复注册检查 - controller
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/frame/gmvc"
)

type User struct {
	gmvc.Controller
}

func (u *User) Index() {
	u.Response.Write("User")
}

func (u *User) Info() {
	u.Response.Write("Info - Uid: ", u.Request.Get("uid"))
}

func (u *User) List() {
	u.Response.Write("List - Page: ", u.Request.Get("page"))
}

func main() {
	s := g.Server()
	s.BindController("/user", new(User))
	s.BindController("/user/{.method}/{uid}", new(User), "Info")
	s.BindController("/user/{.method}/{page}.html", new(User), "List")
	s.BindController("/user/{.method}/{page}.html", new(User), "List")
	s.SetPort(8199)
	s.Run()
}

package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type User struct {
	gmvc.Controller
}

func init() {
	s := g.Server()
	s.BindController("/user", new(User))
	s.BindController("/user/{.method}/{uid}", new(User), "Info")
	s.BindController("/user/{.method}/{page}.html", new(User), "List")
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

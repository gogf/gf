package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type User struct {
}

func (c *User) Index(r *ghttp.Request) {
	r.Response.Write("Index")
}

// 不符合规范，不会被注册
func (c *User) Test(r *ghttp.Request, value interface{}) {
	r.Response.Write("Test")
}

func main() {
	s := g.Server()
	s.BindObject("/user", new(User))
	s.SetPort(8199)
	s.Run()
}

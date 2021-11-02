package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type User struct{}

func (c *User) Test(r *ghttp.Request) {
	r.Response.Write("Test")
}

func main() {
	s := g.Server()
	u := new(User)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/db-{table}/{id}", u, "Test")
	})
	s.SetPort(8199)
	s.Run()
}

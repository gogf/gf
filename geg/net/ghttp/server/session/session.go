package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type Controller struct {
	gmvc.Controller
}

func (c *Controller) Login() {
	c.Session.Id()
	c.Response.Write("这个页面用户填写信息执行登录")
}

func (c *Controller) DoLogin() {
	c.Session.Set("key", "value")
	//c.Response.Header().Set("Set-Cookie", "myid=1B27UGQGCIBP0P70; Path=/; Domain=127.0.0.1; Expires=Wed, 04 Mar 2020 07:12:05 GMT")

	c.Response.RedirectTo("/main")
}

func (c *Controller) Main() {
	c.Response.WriteJson(c.Session.Map())
}

func main() {
	s := g.Server()
	s.BindController("/", new(Controller))
	s.SetPort(8199)
	s.Run()
}

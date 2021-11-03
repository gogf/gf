// 路由重复注册检查 - object
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Object struct{}

func (o *Object) Index(r *ghttp.Request) {
	r.Response.Write("object index")
}

func (o *Object) Show(r *ghttp.Request) {
	r.Response.Write("object show")
}

func main() {
	s := g.Server()
	g.Server().BindObject("/object", new(Object))
	g.Server().BindObject("/object", new(Object))
	s.SetPort(8199)
	s.Run()
}

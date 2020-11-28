package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type GetById struct {
	Id *g.Var `p:"id" v:"required|integer#id不能为空|id必须为整数"`
}

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
		var idInfo *GetById
		if err := r.Parse(&idInfo); err != nil {
			r.Response.Write(err)
		}
		r.Response.Write("ok")
	})
	s.SetPort(8999)
	s.Run()
}

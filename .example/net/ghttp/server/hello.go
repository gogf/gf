package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type GetById struct {
	Id *g.Var `p:"id" v:"required|integer#id不能为空|id必须为整数"`
}

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
<<<<<<< HEAD
		glog.Println(r.Header)
		r.Response.Write("hello world")
=======
		var idInfo *GetById
		if err := r.Parse(&idInfo); err != nil {
			r.Response.Write(err)
		}
		r.Response.Write("ok")
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	})
	s.SetPort(8999)
	s.Run()
}

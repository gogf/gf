package main

import (
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
)

func main() {
	gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(`Hello world`)
	})
	s.Run()
}

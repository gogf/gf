package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	// 该路由规则仅会在GET请求下有效
	s.BindHandler("GET:/{table}/list/{page}.html", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
	})
	// 该路由规则仅会在GET请求及localhost域名下有效
	s.BindHandler("GET:/order/info/{order_id}@localhost", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
	})
	// 该路由规则仅会在DELETE请求下有效
	s.BindHandler("DELETE:/comment/{id}", func(r *ghttp.Request) {
		r.Response.WriteJson(r.Router)
	})
	s.SetPort(8199)
	s.Run()
}

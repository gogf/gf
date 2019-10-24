package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func loadRouter(domain *ghttp.Domain) {
	domain.Group("/", func(g *ghttp.RouterGroup) {
		g.Group("/app", func(gApp *ghttp.RouterGroup) {
			// 该路由规则仅会在GET请求下有效
			gApp.GET("/{table}/list/{page}.html", func(r *ghttp.Request) {
				r.Response.WriteJson(r.Router)
			})
			// 该路由规则仅会在GET请求及localhost域名下有效
			gApp.GET("/order/info/{order_id}", func(r *ghttp.Request) {
				r.Response.WriteJson(r.Router)
			})
			// 该路由规则仅会在DELETE请求下有效
			gApp.DELETE("/comment/{id}", func(r *ghttp.Request) {
				r.Response.WriteJson(r.Router)
			})
		})
		// 该路由规则仅会在GET请求下有效
		g.GET("/{table}/list/{page}.html", func(r *ghttp.Request) {
			r.Response.WriteJson(r.Router)
		})
		// 该路由规则仅会在GET请求及localhost域名下有效
		g.GET("/order/info/{order_id}", func(r *ghttp.Request) {
			r.Response.WriteJson(r.Router)
		})
		// 该路由规则仅会在DELETE请求下有效
		g.DELETE("/comment/{id}", func(r *ghttp.Request) {
			r.Response.WriteJson(r.Router)
		})
	})
}

func main() {
	s := g.Server()

	domain := s.Domain("localhost")
	loadRouter(domain)

	s.SetPort(8199)
	s.Run()
}

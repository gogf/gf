package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
)

func main() {
	type User struct {
		Id   int64
		Name string
	}
	s := g.Server()
	s.SetSessionStorage(gsession.NewStorageRedis(g.Redis()))
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/set", func(r *ghttp.Request) {
			user := &User{
				Id:   1265476890672672808,
				Name: "john",
			}
			if err := r.Session.Set("user", user); err != nil {
				panic(err)
			}
			r.Response.Write("ok")
		})
		group.GET("/get", func(r *ghttp.Request) {
			r.Response.WriteJson(r.Session.Get("user"))
		})
		group.GET("/clear", func(r *ghttp.Request) {
			r.Session.Clear()
		})
	})
	s.SetPort(8199)
	s.Run()
}

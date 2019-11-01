package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gsession"
	"github.com/gogf/gf/os/gtime"
	"time"
)

// SessionSet sets a test session key into the session storage.
func SessionSet(r *ghttp.Request) {
	r.Session.Set("time", gtime.Second())
	r.Response.WriteJson("ok")
}

// SessionGet shows all sessions stored.
func SessionGet(r *ghttp.Request) {
	r.Response.WriteJson(r.Session.Map())
}

func main() {
	s := g.Server()
	s.SetConfigWithMap(g.Map{
		"SessionMaxAge":  3 * time.Second,
		"SessionStorage": gsession.NewStorageRedis(g.Redis()),
	})
	s.BindHandler("/set", SessionSet)
	s.BindHandler("/get", SessionGet)
	s.SetPort(8199)
	s.Run()
}

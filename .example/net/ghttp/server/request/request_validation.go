package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/util/gvalid"
)

type User struct {
	Uid   int    `gvalid:"uid@min:1"`
	Name  string `params:"username"  gvalid:"username @required|length:6,30"`
	Pass1 string `params:"password1" gvalid:"password1@required|password3"`
	Pass2 string `params:"password2" gvalid:"password2@required|password3|same:password1#||两次密码不一致，请重新输入"`
}

func main() {
	s := g.Server()
	s.Group("/", func(rgroup *ghttp.RouterGroup) {
		rgroup.ALL("/user", func(r *ghttp.Request) {
			user := new(User)
			if err := r.GetToStruct(user); err != nil {
				r.Response.WriteJsonExit(g.Map{
					"message": err,
					"errcode": 1,
				})
			}
			if err := gvalid.CheckStruct(user, nil); err != nil {
				r.Response.WriteJsonExit(g.Map{
					"message": err.Maps(),
					"errcode": 1,
				})
			}
			r.Response.WriteJsonExit(g.Map{
				"message": "ok",
				"errcode": 0,
			})
		})
	})
	s.SetPort(8199)
	s.Run()
}

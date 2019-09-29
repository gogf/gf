package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	s := g.Server()
	// 多事件回调示例，事件1
	pattern1 := "/:name/info"
	s.BindHookHandlerByMap(pattern1, map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.SetParam("uid", 1000)
		},
	})
	s.BindHandler(pattern1, func(r *ghttp.Request) {
		r.Response.Write("用户:", r.Get("name"), ", uid:", r.Get("uid"))
	})

	// 多事件回调示例，事件2
	pattern2 := "/{object}/list/{page}.java"
	s.BindHookHandlerByMap(pattern2, map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_OUTPUT: func(r *ghttp.Request) {
			r.Response.SetBuffer([]byte(
				fmt.Sprintf("通过事件修改输出内容, object:%s, page:%s", r.Get("object"), r.GetRouterString("page"))),
			)
		},
	})
	s.BindHandler(pattern2, func(r *ghttp.Request) {
		r.Response.Write(r.Router.Uri)
	})
	s.SetPort(8199)
	s.Run()
}

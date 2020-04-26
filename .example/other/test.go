package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtimer"
)

func main() {
	s := g.Server()
	s.SetDumpRouterMap(false)

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			paramsMap := r.GetRequestMap()
			fmt.Print("打印参数\n", paramsMap)
		})
	})

	addr := "localhost:8199"
	gtimer.SetTimeout(time.Second, func() {
		client := g.Client().SetHeader("Content-Type", "application/x-www-form-urlencoded")
		client.PostContent(
			fmt.Sprintf("http://%s", addr),
			"time_end2020-04-18 16:11:58&returnmsg=Success&attach=",
		)

		fmt.Print("\n")
		client.PostContent(
			fmt.Sprintf("http://%s", addr),
			"returnmsg=Success&attach=",
		)
	})

	s.SetAddr(addr)
	s.Run()
}

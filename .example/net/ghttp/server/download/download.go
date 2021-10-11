package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	s.BindHandler("/download", func(r *ghttp.Request) {
		r.Response.Header().Set("Content-Type", "text/html;charset=utf-8")
		r.Response.Header().Set("Content-type", "application/force-download")
		r.Response.Header().Set("Content-Type", "application/octet-stream")
		r.Response.Header().Set("Accept-Ranges", "bytes")
		r.Response.Header().Set("Content-Disposition", "attachment;filename=\"下载文件名称.txt\"")
		r.Response.ServeFile("text.txt")
	})
	s.SetPort(8199)
	s.Run()
}

package main

import "github.com/gogf/gf/v2/frame/g"

// 静态文件服务器基本使用
func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.SetServerRoot("/Users/john/Downloads")
	//s.AddSearchPath("/Users/john/Documents")
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(8199)
	s.Run()
}

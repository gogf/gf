package main

import "gitee.com/johng/gf/g"

// 静态文件服务器基本使用
func main() {
    s := g.Server()
    s.SetIndexFolder(true)
    s.SetServerRoot("/Users/john/Temp")
    s.AddSearchPath("/Users/john/Documents")
    s.SetPort(8199)
    s.Run()
}

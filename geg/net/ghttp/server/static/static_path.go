package main

import "gitee.com/johng/gf/g"

// 静态文件服务器，支持自定义静态目录映射
func main() {
    s := g.Server()
    s.SetIndexFolder(true)
    s.SetServerRoot("/Users/john/Temp")
    s.AddSearchPath("/Users/john/Documents")
    s.AddStaticPath("/my-doc", "/Users/john/Documents")
    s.SetPort(8199)
    s.Run()
}

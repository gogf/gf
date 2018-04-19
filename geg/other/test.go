package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.SetServerRoot("/home/john/Documents")
    s.SetIndexFolder(true)
    s.SetPort(8199)
    s.Run()
}
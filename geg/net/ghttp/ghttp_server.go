package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.SetAddr(":8199")
    s.SetIndexFolder(true)
    s.SetServerRoot("/tmp")
    s.Run()
}

package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/demo"
)

func main() {
    ghttp.GetServer().SetPort(8199)
    ghttp.GetServer().Run()
}

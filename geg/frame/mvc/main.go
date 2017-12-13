package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/user"

    "gitee.com/johng/gf/g/frame/gconfig"
)

func main() {
    gconfig.Set("johng.gf.mvc.view.path", "/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/mvc/view")
    ghttp.GetServer("johng").SetAddr(":8199")
    ghttp.GetServer("johng").Run()
}

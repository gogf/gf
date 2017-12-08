package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/user"

)

func main() {
    ghttp.GetServer("johng.cn").SetAddr(":8199")
    ghttp.GetServer("johng.cn").Run()
}

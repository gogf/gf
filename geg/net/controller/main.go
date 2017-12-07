package main

import (
    _ "gitee.com/johng/gf/geg/net/controller/controller"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    ghttp.GetServer("johng.cn").SetAddr(":8199")
    ghttp.GetServer("johng.cn").Run()
}

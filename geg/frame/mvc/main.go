package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/user"
)

func main() {
    g.HttpServer().SetPort(8199)
    g.HttpServer().Run()
}

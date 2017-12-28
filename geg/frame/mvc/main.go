package main

import (
    "gitee.com/johng/gf/g/frame/ginstance"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/user"
)

func main() {
    ginstance.Server().SetPort(8199)
    ginstance.Server().Run()
}

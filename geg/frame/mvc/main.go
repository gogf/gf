package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/demo"
)

func main() {
    g.HTTPServer().SetPort(8199)
    g.HTTPServer().Run()
}

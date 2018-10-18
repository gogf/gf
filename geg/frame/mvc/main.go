package main

import (
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/demo"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/stats"
)

func main() {

    //g.Server().SetDumpRouteMap(false)
    g.Server().SetPort(8199)
    g.Server().Run()

}

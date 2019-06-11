package main

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g"
    _ "gitee.com/johng/gf/geg/frame/mvc/controller/demo"
)

func main() {
    g.Server().SetPort(8199)
    g.Server().Run()
=======
	"github.com/gogf/gf/g"
	_ "github.com/gogf/gf/geg/frame/mvc/controller/demo"
	_ "github.com/gogf/gf/geg/frame/mvc/controller/stats"
)

func main() {

	//g.Server().SetDumpRouteMap(false)
	g.Server().SetPort(8199)
	g.Server().Run()

>>>>>>> upstream/master
}

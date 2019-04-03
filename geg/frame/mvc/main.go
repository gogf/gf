package main

import (
	"github.com/gogf/gf/g"
	_ "github.com/gogf/gf/geg/frame/mvc/controller/demo"
	_ "github.com/gogf/gf/geg/frame/mvc/controller/stats"
)

func main() {

	//g.Server().SetDumpRouteMap(false)
	g.Server().SetPort(8199)
	g.Server().Run()

}

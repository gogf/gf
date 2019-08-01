package main

import (
	_ "github.com/gogf/gf/.example/frame/mvc/controller/demo"
	_ "github.com/gogf/gf/.example/frame/mvc/controller/stats"
	"github.com/gogf/gf/frame/g"
)

func main() {

	//g.Server().SetDumpRouteMap(false)
	g.Server().SetPort(8199)
	g.Server().Run()

}

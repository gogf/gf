package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	g.TryCatch(func() {
		glog.Printfln("hello")
		g.Throw("exception")
		glog.Printfln("world")
	}, func(exception interface{}) {
		glog.Error(exception)
	})
}

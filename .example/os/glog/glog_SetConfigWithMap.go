package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	err := g.Log().SetConfigWithMap(g.Map{
		"prefix": "[TEST]",
	})
	if err != nil {
		panic(err)
	}
	glog.Info(1)
}

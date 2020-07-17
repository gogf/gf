package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/glog"
)

func main() {
	err := glog.SetConfigWithMap(g.Map{
		"prefix": "[TEST]",
	})
	if err != nil {
		panic(err)
	}
	glog.Info(1)
}

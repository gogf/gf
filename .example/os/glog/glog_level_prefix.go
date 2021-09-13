package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	g.Log().SetLevelPrefix(glog.LEVEL_DEBU, "debug")
	g.Log().Debug("test")
}

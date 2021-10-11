package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	g.Log().SetLevelPrefix(glog.LEVEL_DEBU, "debug")
	g.Log().Debug("test")
}

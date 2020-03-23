package main

import (
	"github.com/gogf/gf/os/glog"
)

func main() {
	l := glog.New()
	l.Info("info1")
	l.SetLevelStr("notice")
	l.Info("info2")
}

package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	glog.SetDebug(false)
	glog.Warning(1)
	glog.SetDebug(true)
	glog.Warning(1)
}

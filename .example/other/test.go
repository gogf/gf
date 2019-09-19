package main

import (
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/glog"
)

func main() {
	glog.SetFlags(glog.F_TIME_DATE | glog.F_TIME_TIME | glog.F_FILE_SHORT)
	glog.Debug("dd")
	glog.Println("timeout", gcmd.GetOpt("timeout"))
}

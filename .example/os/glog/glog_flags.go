package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	g.Log().SetFlags(glog.F_TIME_TIME | glog.F_FILE_SHORT)
	g.Log().Println("time and short line number")
	g.Log().SetFlags(glog.F_TIME_MILLI | glog.F_FILE_LONG)
	g.Log().Println("time with millisecond and long line number")
	g.Log().SetFlags(glog.F_TIME_STD | glog.F_FILE_LONG)
	g.Log().Println("standard time format and long line number")
}

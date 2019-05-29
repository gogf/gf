package main

import (
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gtime"
)

func main() {
	Time := gtime.Now().AddDate(0, -1, 0).Format("Y-m")
	glog.Debug(Time)
	Time = gtime.Now().AddDate(0, -2, 0).Format("Y-m")
	glog.Debug(Time)
	Time = gtime.Now().AddDate(0, -3, 0).Format("Y-m")
	glog.Debug(Time)
	Time = gtime.Now().AddDate(0, -4, 0).Format("Y-m")
	glog.Debug(Time)
}
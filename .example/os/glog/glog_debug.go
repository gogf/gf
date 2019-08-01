package main

import (
	"time"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
)

func main() {
	gtimer.SetTimeout(3*time.Second, func() {
		glog.SetDebug(false)
	})
	for {
		glog.Debug(gtime.Datetime())
		time.Sleep(time.Second)
	}
}

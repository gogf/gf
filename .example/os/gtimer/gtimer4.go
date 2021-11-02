package main

import (
	"time"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	interval := time.Second
	gtimer.AddTimes(interval, 2, func() {
		glog.Print("doing1")
	})
	gtimer.AddTimes(interval, 2, func() {
		glog.Print("doing2")
	})

	select {}
}

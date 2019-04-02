package main

import (
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gtimer"
	"time"
)

func main() {
	interval := time.Second
	gtimer.AddSingleton(interval, func() {
		glog.Println("doing")
		time.Sleep(5 * time.Second)
	})

	select {}
}

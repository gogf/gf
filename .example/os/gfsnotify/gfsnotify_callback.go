package main

import (
	"time"

	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	c1, err := gfsnotify.Add("/home/john/temp/log", func(event *gfsnotify.Event) {
		glog.Println("callback1")
	})
	if err != nil {
		panic(err)
	}
	c2, err := gfsnotify.Add("/home/john/temp/log", func(event *gfsnotify.Event) {
		glog.Println("callback2")
	})
	if err != nil {
		panic(err)
	}
	// 5秒后移除c1的回调函数注册，仅剩c2
	gtimer.SetTimeout(5*time.Second, func() {
		gfsnotify.RemoveCallback(c1.Id)
		glog.Println("remove callback c1")
	})
	// 10秒后移除c2的回调函数注册，所有的回调都移除，不再有任何打印信息输出
	gtimer.SetTimeout(10*time.Second, func() {
		gfsnotify.RemoveCallback(c2.Id)
		glog.Println("remove callback c2")
	})

	select {}

}

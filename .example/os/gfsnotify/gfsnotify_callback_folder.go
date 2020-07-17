package main

import (
	"time"

	"github.com/jin502437344/gf/os/gfsnotify"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gtimer"
)

func main() {
	callback, err := gfsnotify.Add("/home/john/temp", func(event *gfsnotify.Event) {
		glog.Println("callback")
	})
	if err != nil {
		panic(err)
	}

	// 在此期间创建文件、目录、修改文件、删除文件

	// 20秒后移除回调函数注册，所有的回调都移除，不再有任何打印信息输出
	gtimer.SetTimeout(20*time.Second, func() {
		gfsnotify.RemoveCallback(callback.Id)
		glog.Println("remove callback")
	})

	select {}
}

package main

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtimer"
)

func main() {
	var (
		ctx = context.Background()
	)
	callback, err := gfsnotify.Add("/home/john/temp", func(event *gfsnotify.Event) {
		glog.Print(ctx, "callback")
	})
	if err != nil {
		panic(err)
	}

	// 在此期间创建文件、目录、修改文件、删除文件

	// 20秒后移除回调函数注册，所有的回调都移除，不再有任何打印信息输出
	gtimer.SetTimeout(ctx, 20*time.Second, func(ctx context.Context) {
		gfsnotify.RemoveCallback(callback.Id)
		glog.Print(ctx, "remove callback")
	})

	select {}
}

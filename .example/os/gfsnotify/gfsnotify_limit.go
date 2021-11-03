package main

import (
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/glog"
)

// 对同一个文件多次Add是否超过系统inotify限制
func main() {
	path := "/Users/john/temp/log"
	for i := 0; i < 9999999; i++ {
		_, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
			glog.Print(event)
		})
		if err != nil {
			glog.Fatal(err)
		}
	}
}

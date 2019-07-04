package main

import (
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	path := `D:\temp`
	//path := "/Users/john/Temp/test"
	_, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
		glog.Println(event)
	}, true)
	if err != nil {
		glog.Fatal(err)
	} else {
		select {}
	}
}

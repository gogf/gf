package main

import (
<<<<<<< HEAD
    "log"
    "gitee.com/johng/gf/g/os/gfsnotify"
)

func main() {
    err := gfsnotify.Add("./temp.txt", func(event *gfsnotify.Event) {
        if event.IsCreate() {
            log.Println("创建文件 : ", event.Path)
        }
        if event.IsWrite() {
            log.Println("写入文件 : ", event.Path)
        }
        if event.IsRemove() {
            log.Println("删除文件 : ", event.Path)
        }
        if event.IsRename() {
            log.Println("重命名文件 : ", event.Path)
        }
        if event.IsChmod() {
            log.Println("修改权限 : ", event.Path)
        }
    })
    if err != nil {
        log.Fatalln(err)
    } else {
        select {}
    }
}
=======
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/os/glog"
)

func main() {
	//path := "D:\\Workspace\\Go\\GOPATH\\src\\gitee.com\\johng\\gf\\geg\\other\\test.go"
	path := "/Users/john/Temp/test"
	_, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
		glog.Println(event)
	}, true)
	if err != nil {
		glog.Fatal(err)
	} else {
		select {}
	}
}
>>>>>>> upstream/master

package main

import (
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    // /home/john/temp 是一个目录，当然也可以指定文件
    path := "/home/john/temp"
    _, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
        if event.IsCreate() {
            glog.Println("创建文件 : ", event.Path)
        }
        if event.IsWrite() {
            glog.Println("写入文件 : ", event.Path)
        }
        if event.IsRemove() {
            glog.Println("删除文件 : ", event.Path)
        }
        if event.IsRename() {
            glog.Println("重命名文件 : ", event.Path)
        }
        if event.IsChmod() {
            glog.Println("修改权限 : ", event.Path)
        }
        glog.Println(event)
    })

    // 移除对该path的监听
    gfsnotify.Remove(path)

    if err != nil {
        glog.Fatalln(err)
    } else {
        select {}
    }
}
package main

import (
    "log"
    "gitee.com/johng/gf/g/os/gfsnotify"
)

func main() {
    err := gfsnotify.Add("/home/john/Documents/temp", func(event *gfsnotify.Event) {
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
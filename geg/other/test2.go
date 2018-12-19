package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    path := "/Users/john/Temp"
    _, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
        fmt.Println(event)
        if event.IsWrite() {
            glog.Println("写入文件 : ", event.Path)
            fmt.Printf("%s\n", gfile.GetContents(event.Path))
        }
    })
    if err != nil {
        glog.Fatal(err)
    } else {
        select {}
    }

}
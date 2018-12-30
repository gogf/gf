package main

import (
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/os/glog"
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
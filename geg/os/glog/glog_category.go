package main

import (
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    glog.SetPath("/tmp/")
    glog.Cat("test1").Cat("test2").Println("test")
    glog.Println("test")
}



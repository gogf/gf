package main

import (
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main() {
    glog.SetPath("/tmp/test-logs")
    for {
        glog.Println("1")
        time.Sleep(time.Second)
    }
}

